package eveConsumer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/v1"
	"github.com/garyburd/redigo/redis"

	"golang.org/x/oauth2"
)

func init() {
	addConsumer("contactSync", contactSyncConsumer, "EVEDATA_contactSyncQueue")
	addTrigger("contactSync", contactSyncTrigger)
}

// Perform contact sync for wardecs
func contactSyncTrigger(c *EVEConsumer) (bool, error) {

	// Do quick maintenence to prevent errors.
	err := models.MaintContactSync()
	if err != nil {
		return false, err
	}

	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT source, group_concat(destination)
			FROM evedata.contactSyncs S  
            INNER JOIN evedata.crestTokens T ON T.tokenCharacterID = destination
            WHERE lastStatus NOT LIKE "%400 Bad Request%"
		    GROUP BY source
            HAVING max(nextSync) < UTC_TIMESTAMP();`)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	} else if err == sql.ErrNoRows { // Shut up warnings
		return false, nil
	}
	defer rows.Close()

	r := c.ctx.Cache.Get()
	defer r.Close()

	// Loop updatable characters
	for rows.Next() {
		var (
			source int64  // Source char
			dest   string // List of destination chars
		)

		err = rows.Scan(&source, &dest)
		if err != nil {
			log.Printf("Contact Sync: Failed scan: %v", err)
			continue
		}

		_, err = r.Do("SADD", "EVEDATA_contactSyncQueue", fmt.Sprintf("%d:%s", source, dest))
		if err != nil {
			log.Printf("Contact Sync: Failed scan: %v", err)
			continue
		}
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
	destinations := strings.Split(dest[1], ",")
	source, err := strconv.ParseInt(dest[0], 10, 64)
	if err != nil {
		return false, err
	}

	// get the source character information
	char, _, err := c.ctx.ESI.V4.CharacterApi.GetCharactersCharacterId((int32)(source), nil)
	if err != nil {
		return false, err
	}

	corp, _, err := c.ctx.ESI.V3.CorporationApi.GetCorporationsCorporationId(char.CorporationId, nil)
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
		token *oauth2.TokenSource
		cid   int64
	}
	tokens := make(map[int64]characterToken)

	// Get the tokens for our destinations
	for _, cidS := range destinations {
		cid, _ := strconv.ParseInt(cidS, 10, 64)
		a, err := c.getToken(source, cid)
		if err != nil {
			return false, err
		}
		// Save the token.
		tokens[cid] = characterToken{token: &a, cid: cid}
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

	// Loop through all the destinations
	for _, token := range tokens {
		// authentication token context for destination char
		auth := context.WithValue(context.TODO(), goesi.ContextOAuth2, *token.token)
		var (
			contacts []goesiv1.GetCharactersCharacterIdContacts200Ok
			r        *http.Response
			err      error
		)

		// Default to OK
		tokenSuccess(source, token.cid, 200, "OK")

		// Get current contacts
		for i := 1; ; i++ {
			var con []goesiv1.GetCharactersCharacterIdContacts200Ok
			con, r, err = c.ctx.ESI.V1.ContactsApi.GetCharactersCharacterIdContacts(auth, (int32)(token.cid), map[string]interface{}{"page": (int32)(i)})
			if err != nil || r.StatusCode != 200 {
				tokenError(source, token.cid, r, err)
				return false, err
			}
			if len(con) == 0 {
				break
			}
			contacts = append(contacts, con...)
		}

		// Update cache time.
		if r != nil {
			contactSync := &models.ContactSync{Source: source, Destination: token.cid}
			err := contactSync.Updated(goesi.CacheExpires(r))
			if err != nil {
				return false, err
			}
		}

		var erase []int32
		var active []int32
		var pending []int32
		var pendingMove []int32
		var activeMove []int32

		activeCheck := make(map[int32]bool)
		pendingCheck := make(map[int32]bool)

		for _, war := range activeWars {
			activeCheck[(int32)(war.ID)] = true
		}
		for _, war := range pendingWars {
			pendingCheck[(int32)(war.ID)] = true
		}

		// Loop through all current contacts
		for _, contact := range contacts {
			// skip anything > -0.4
			if contact.Standing > -0.4 {
				continue
			}

			_, pend := pendingCheck[contact.ContactId]
			_, act := activeCheck[contact.ContactId]

			// Is this existing contact in the pending list
			if !pend {
				// Is this existing contact in the active list
				if !act { // Not in either list. delete it.
					erase = append(erase, (int32)(contact.ContactId))
				} else if act && contact.Standing > -10.0 { // in active list but wrong standing
					// Take it out of the pending list and put into active move.
					delete(activeCheck, contact.ContactId)
					activeMove = append(activeMove, (int32)(contact.ContactId))
				} else if act && contact.Standing == -10.0 { // Contact correct, do nothing.
					delete(activeCheck, contact.ContactId)
				}
			} else if pend && contact.Standing != -5.0 { // in pending list, but wrong standing
				delete(pendingCheck, contact.ContactId)
				pendingMove = append(pendingMove, (int32)(contact.ContactId))
			} else if pend && contact.Standing == -5.0 { // Contact correct, do nothing.
				delete(pendingCheck, contact.ContactId)
			}
		}

		for con, _ := range activeCheck {
			active = append(active, con)
		}

		for con, _ := range pendingCheck {
			pending = append(pending, con)
		}

		if len(erase) > 0 {
			for start := 0; start < len(erase); start = start + 20 {
				end := min(start+20, len(erase))
				r, err = c.ctx.ESI.V1.ContactsApi.DeleteCharactersCharacterIdContacts(auth, (int32)(token.cid), erase[start:end], nil)
				if err != nil {
					tokenError(source, token.cid, r, err)
					return false, err
				}
			}
		}
		if len(active) > 0 {
			for start := 0; start < len(active); start = start + 100 {
				end := min(start+100, len(active))
				_, r, err = c.ctx.ESI.V1.ContactsApi.PostCharactersCharacterIdContacts(auth, (int32)(token.cid), active[start:end], -10, nil)

				if err != nil {
					tokenError(source, token.cid, r, err)
					return false, err
				}
			}
		}
		if len(pending) > 0 {
			for start := 0; start < len(pending); start = start + 100 {
				end := min(start+100, len(pending))
				_, r, err = c.ctx.ESI.V1.ContactsApi.PostCharactersCharacterIdContacts(auth, (int32)(token.cid), pending[start:end], -5, nil)
				if err != nil {
					tokenError(source, token.cid, r, err)
					return false, err
				}
			}
		}
		if len(activeMove) > 0 {
			for start := 0; start < len(activeMove); start = start + 20 {
				end := min(start+20, len(activeMove))
				r, err = c.ctx.ESI.V1.ContactsApi.PutCharactersCharacterIdContacts(auth, (int32)(token.cid), activeMove[start:end], -10, nil)
				if err != nil {
					tokenError(source, token.cid, r, err)
					return false, err
				}
			}
		}
		if len(pendingMove) > 0 {
			for start := 0; start < len(pendingMove); start = start + 20 {
				end := min(start+20, len(pendingMove))
				r, err = c.ctx.ESI.V1.ContactsApi.PutCharactersCharacterIdContacts(auth, (int32)(token.cid), pendingMove[start:end], -5, nil)
				if err != nil {
					tokenError(source, token.cid, r, err)
					return false, err
				}
			}
		}
	}
	return true, err
}
