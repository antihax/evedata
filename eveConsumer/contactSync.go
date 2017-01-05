package eveConsumer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/antihax/evedata/esi"
	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"

	"golang.org/x/oauth2"
)

// Perform contact sync for wardecs
func (c *EVEConsumer) contactSync() {
	r := c.ctx.Cache.Get()
	defer r.Close()

	log.Printf("Running Contact Sync\n")

	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT source, group_concat(destination)
			FROM contactSyncs S  
            INNER JOIN crestTokens T ON T.tokenCharacterID = destination
            WHERE lastStatus NOT LIKE "%Invalid refresh token%"
		    GROUP BY source
            HAVING max(nextSync) < UTC_TIMESTAMP();`)
	if err != nil {
		log.Printf("Contact Sync: Failed query: %v", err)
		return
	}

	defer rows.Close()

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
}

func (c *EVEConsumer) contactSyncCheckQueue(r redis.Conn) error {

	ret, err := r.Do("SPOP", "EVEDATA_contactSyncQueue")
	if err != nil {
		return err
	} else if ret == nil {
		return nil
	}

	v, err := redis.String(ret, err)
	if err != nil {
		return err
	}

	// Split off characters into an array
	dest := strings.Split(v, ":")
	destinations := strings.Split(dest[1], ",")
	source, err := strconv.ParseInt(dest[0], 10, 64)
	if err != nil {
		return err
	}
	// get the source character information
	char, err := c.ctx.EVE.CharacterInfoXML(source)
	if err != nil {
		return err
	}

	// Find the Entity ID to search for wars.
	var searchID int64
	if char.AllianceID > 0 {
		searchID = char.AllianceID
	} else {
		searchID = char.CharacterID
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
			log.Printf("Contact Sync: Failed getClient %d %d %v", source, cid, err)
		}
		// Save the token.
		tokens[cid] = characterToken{token: &a, cid: cid}
	}

	// Active Wars
	activeWars, err := models.GetActiveWarsByID(searchID)
	if err != nil {
		log.Printf("Contact Sync: Failed Getting Active Wars: %v", err)
	}

	// Pending Wars
	pendingWars, err := models.GetPendingWarsByID(searchID)
	if err != nil {
		log.Printf("Contact Sync: Failed Getting Pending Wars: %v", err)
	}

	// Make a list of contacts to add.
	var pending []int32
	var active []int32
	pendingToAdd := make(map[int32]int32)
	activeToAdd := make(map[int32]int32)
	for _, war := range activeWars {
		activeToAdd[(int32)(war.ID)] = (int32)(war.ID)
		active = append(active, (int32)(war.ID))
	}
	for _, war := range pendingWars {
		pendingToAdd[(int32)(war.ID)] = (int32)(war.ID)
		pending = append(pending, (int32)(war.ID))
	}

	// Loop through all the destinations
	for _, token := range tokens {
		// authentication token context for destination char
		auth := context.WithValue(context.TODO(), esi.ContextOAuth2, *token.token)
		var (
			contacts []esi.GetCharactersCharacterIdContacts200Ok
			r        *http.Response
			err      error
		)

		syncSuccess(source, token.cid, 200, "OK")

		// Get current contacts
		for i := 1; ; i++ {
			var con []esi.GetCharactersCharacterIdContacts200Ok
			con, r, err = c.ctx.ESI.ContactsApi.GetCharactersCharacterIdContacts(auth, (int32)(token.cid), map[string]interface{}{"page": (int32)(i)})
			if err != nil {
				syncError(source, token.cid, r, err)
				break
			}
			if len(con) == 0 {
				break
			}
			contacts = append(contacts, con...)
		}

		// Update cache time.
		if r != nil {
			contactSync := &models.ContactSync{Source: source, Destination: token.cid}
			contactSync.Updated(esi.CacheExpires(r))
		}

		var erase []int32

		// Loop through all current contacts
		for _, contact := range contacts {
			// skip anything > -0.4
			if contact.Standing > -0.4 {
				continue
			}
			erase = append(erase, (int32)(contact.ContactId))
			/*if _, ok := pendingToAdd[contact.ContactId]; !ok {
				if _, ok := activeToAdd[contact.ContactId]; !ok {
					erase = append(erase, (int32)(contact.ContactId))
				}
			}*/
		}
		if len(erase) > 0 {
			for start := 0; start < len(erase); start = start + 20 {
				end := min(start+20, len(erase))
				r, err = c.ctx.ESI.ContactsApi.DeleteCharactersCharacterIdContacts(auth, (int32)(token.cid), erase[start:end], nil)
				if err != nil {
					syncError(source, token.cid, r, err)
					break
				}
			}
		}
		if len(active) > 0 {
			for start := 0; start < len(active); start = start + 100 {
				end := min(start+100, len(active))
				_, r, err = c.ctx.ESI.ContactsApi.PostCharactersCharacterIdContacts(auth, (int32)(token.cid), -10, active[start:end], nil)
				if err != nil {
					syncError(source, token.cid, r, err)
					break
				}
			}
		}
		if len(pending) > 0 {
			for start := 0; start < len(pending); start = start + 100 {
				end := min(start+100, len(pending))
				_, r, err = c.ctx.ESI.ContactsApi.PostCharactersCharacterIdContacts(auth, (int32)(token.cid), -5, pending[start:end], nil)
				if err != nil {
					syncError(source, token.cid, r, err)
					break
				}
			}
		}
	}

	return err
}
