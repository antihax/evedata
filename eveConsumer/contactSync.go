package eveConsumer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("contactSync", contactSyncConsumer, "EVEDATA_contactSyncQueue")
	addTrigger("contactSync", contactSyncTrigger)
}

// Perform contact sync for wardecs
func contactSyncTrigger(c *EVEConsumer) (bool, error) {
	ecc, err := models.GetExpiredContactSyncs()
	if err != nil {
		return false, err
	}

	r := c.ctx.Cache.Get()
	defer r.Close()

	for _, cc := range ecc {
		err = r.Send("SADD", "EVEDATA_contactSyncQueue",
			fmt.Sprintf("%d:%d:%s", cc.CharacterID, cc.Source, cc.Destinations))
		if err != nil {
			log.Printf("Contact Sync: Send Failed: %v", err)
			continue
		}
	}

	err = r.Flush()
	if err != nil {
		log.Printf("Contact Sync: Flush Failed: %v", err)
	}

	return true, err
}

func contactSyncConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
	ret, err := r.Do("SPOP", "EVEDATA_contactSyncQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}

	v, err := redis.String(ret, err)
	if err != nil {
		return false, err
	}

	// Split off characters into an array
	dest := strings.Split(v, ":")
	destinations := strings.Split(dest[2], ",")
	characterID, err := strconv.ParseInt(dest[0], 10, 32)
	if err != nil {
		return false, err
	}

	source, err := strconv.ParseInt(dest[1], 10, 32)
	if err != nil {
		return false, err
	}

	char, _, err := c.ctx.ESI.ESI.CharacterApi.GetCharactersCharacterId(nil, int32(source), nil)
	if err != nil {
		return false, err
	}

	corp, _, err := c.ctx.ESI.ESI.CorporationApi.GetCorporationsCorporationId(nil, char.CorporationId, nil)
	if err != nil {
		return false, err
	}

	// Find the Entity ID to search for wars.
	var searchID int32
	if corp.AllianceId > 0 {
		searchID = corp.AllianceId
	} else {
		searchID = char.CorporationId
	}

	// Map of tokens
	type characterToken struct {
		token *goesi.CRESTTokenSource
		cid   int32
	}
	tokens := make(map[int64]characterToken)

	// Get the tokens for our destinations
	for _, cidS := range destinations {
		cid, _ := strconv.ParseInt(cidS, 10, 64)
		a, err := c.ctx.TokenStore.GetTokenSource(int32(characterID), int32(cid))
		if err != nil {
			return false, err
		}
		// Save the token.
		tokens[cid] = characterToken{token: &a, cid: int32(cid)}
	}

	// Active Wars
	activeWars, err := models.GetActiveWarsByID((int64)(searchID))
	if err != nil {
		return false, err
	}

	// Pending Wars
	pendingWars, err := models.GetPendingWarsByID((int64)(searchID))
	if err != nil {
		return false, err
	}

	// Faction Wars
	var factionWars []models.FactionWarEntities
	if corp.Faction != "" {
		factionWars, err = models.GetFactionWarEntitiesForID(models.FactionsByName[corp.Faction])
		if err != nil {
			return false, err
		}
	}

	// Loop through all the destinations
	for _, token := range tokens {
		// authentication token context for destination char
		auth := context.WithValue(context.TODO(), goesi.ContextOAuth2, *token.token)

		contacts, err := c.getContacts(auth, (int32)(token.cid))
		if err != nil {
			return false, err
		}

		// Update cache time.
		contactSync := &models.ContactSync{Source: int32(source), Destination: token.cid}
		err = contactSync.Updated(time.Now().UTC().Add(time.Second * 300))
		if err != nil {
			return false, err
		}

		var erase []int32
		var active []int32
		var pending []int32
		var pendingMove []int32
		var activeMove []int32
		var untouchableContacts int

		// Figure out how many contacts they have outside of ours
		for _, contact := range contacts {
			if contact.Standing > -0.4 {
				untouchableContacts++
			}
		}

		// Faction wars can get over the 1024 contact limit so we need to trim
		// real wars will take precedence.
		trim := len(activeWars) + len(pendingWars)

		activeCheck := make(map[int32]bool)
		pendingCheck := make(map[int32]bool)

		// Build a map of active wars
		for _, war := range activeWars {
			activeCheck[(int32)(war.ID)] = true
		}

		// Add faction wars to the active list
		maxFactionWarLength := min(980-trim-untouchableContacts, len(factionWars))
		for _, war := range factionWars[:maxFactionWarLength] {
			activeCheck[(int32)(war.ID)] = true
		}

		// Build a map of pending wars
		for _, war := range pendingWars {
			pendingCheck[(int32)(war.ID)] = true
		}

		// Loop through all current contacts and figure out needed moves
		for _, contact := range contacts {
			// skip anything > -0.4
			if contact.Standing > -0.4 {
				continue
			}

			pend := pendingCheck[contact.ContactId]
			act := activeCheck[contact.ContactId]

			// Is this existing contact in the active list
			if !act {
				// Is this existing contact in the pending list
				if !pend { // Not in either list. delete it.
					erase = append(erase, (int32)(contact.ContactId))
				} else if pend && contact.Standing > -5.0 { // in pending list but wrong standing
					// Take it out of the active list and put into pending move.
					delete(pendingCheck, contact.ContactId)
					pendingMove = append(pendingMove, (int32)(contact.ContactId))
				} else if pend && contact.Standing == -5.0 { // Contact correct, do nothing.
					delete(pendingCheck, contact.ContactId)
				}
			} else if act && contact.Standing != -10.0 { // in active list, but wrong standing
				delete(activeCheck, contact.ContactId)
				activeMove = append(activeMove, (int32)(contact.ContactId))
			} else if act && contact.Standing == -10.0 { // Contact correct, do nothing.
				delete(activeCheck, contact.ContactId)
			}
		}

		// Build a list of active wars to add
		for con := range activeCheck {
			active = append(active, con)
		}

		// Build a list of pending wars to add
		for con := range pendingCheck {
			pending = append(pending, con)
		}

		// Erase contacts which have no wars.
		if len(erase) > 0 {
			for start := 0; start < len(erase); start = start + 20 {
				end := min(start+20, len(erase))
				if err := c.deleteContacts(auth, (int32)(token.cid), erase[start:end]); err != nil {
					return false, err
				}
			}
		}

		// Add contacts for active wars
		if len(active) > 0 {
			for start := 0; start < len(active); start = start + 100 {
				end := min(start+100, len(active))
				if err := c.addContacts(auth, (int32)(token.cid), active[start:end], -10); err != nil {
					return false, err
				}
			}
		}

		// Add contacts for pending wars
		if len(pending) > 0 {
			for start := 0; start < len(pending); start = start + 100 {
				end := min(start+100, len(pending))
				if err := c.addContacts(auth, (int32)(token.cid), pending[start:end], -5); err != nil {
					return false, err
				}
			}
		}

		// Move contacts to active wars
		if len(activeMove) > 0 {
			for start := 0; start < len(activeMove); start = start + 20 {
				end := min(start+20, len(activeMove))
				if err := c.updateContacts(auth, (int32)(token.cid), activeMove[start:end], -10); err != nil {
					return false, err
				}
			}
		}

		// Move contacts to pending wars
		if len(pendingMove) > 0 {
			for start := 0; start < len(pendingMove); start = start + 20 {
				end := min(start+20, len(pendingMove))
				if err := c.updateContacts(auth, (int32)(token.cid), pendingMove[start:end], -5); err != nil {
					return false, err
				}
			}
		}

		// set success
		tokenSuccess(int32(source), token.cid, 200, "OK")
	}
	return true, err
}
