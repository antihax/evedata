package eveConsumer

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/antihax/evedata/esi"
	"github.com/antihax/evedata/eveapi"
	"github.com/antihax/evedata/models"

	"net/http/httputil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Perform contact sync for wardecs
func (c *EVEConsumer) contactSync() {
	log.Printf("Running Contact Sync\n")
	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT source, group_concat(destination)
			FROM contactSyncs S  
            INNER JOIN crestTokens T ON T.tokenCharacterID = destination
		    GROUP BY source
            HAVING max(nextSync) < UTC_TIMESTAMP()`)
	if err != nil {
		log.Printf("EVEConsumer: Failed query: %v", err)
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
			log.Printf("EVEConsumer: Failed scan: %v", err)
			continue
		}

		// Split off characters into an array
		destinations := strings.Split(dest, ",")

		// get the source character information
		char, err := c.ctx.EVE.CharacterInfoXML(source)
		if err != nil {
			log.Printf("EVEConsumer: Failed getting character info %v", err)
			continue
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
				log.Printf("EVEConsumer: Failed getClient %d %d %v", source, cid, err)
				continue
			}
			// Save the token.
			tokens[cid] = characterToken{token: &a, cid: cid}
		}

		// Active Wars
		activeWars, err := models.GetActiveWarsByID(searchID)
		if err != nil {
			log.Printf("EVEConsumer: Failed Getting Active Wars: %v", err)
			continue
		}

		// Pending Wars
		pendingWars, err := models.GetPendingWarsByID(searchID)
		if err != nil {
			log.Printf("EVEConsumer: Failed Getting Pending Wars: %v", err)
			continue
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
				for start := 0; start < len(active); start = start + 20 {
					end := min(start+20, len(active))
					_, r, err = c.ctx.ESI.ContactsApi.PostCharactersCharacterIdContacts(auth, (int32)(token.cid), -10, active[start:end], nil)
					if err != nil {
						syncError(source, token.cid, r, err)
						break
					}
				}
			}
			if len(pending) > 0 {
				for start := 0; start < len(pending); start = start + 20 {
					end := min(start+20, len(pending))
					_, r, err = c.ctx.ESI.ContactsApi.PostCharactersCharacterIdContacts(auth, (int32)(token.cid), -5, pending[start:end], nil)
					if err != nil {
						syncError(source, token.cid, r, err)
						break
					}
				}
			}
		}
	}
}

// Obtain an authenticated client from a stored access/refresh token.
func (c *EVEConsumer) getToken(characterID int64, tokenCharacterID int64) (oauth2.TokenSource, error) {
	tok := models.CRESTToken{}
	if err := c.ctx.Db.QueryRowx(
		`SELECT expiry, tokenType, accessToken, refreshToken, tokenCharacterID, characterID
			FROM crestTokens
			WHERE characterID = ? AND tokenCharacterID = ?
			LIMIT 1`,
		characterID, tokenCharacterID).StructScan(&tok); err != nil {

		return nil, err
	}

	token := &eveapi.CRESTToken{Expiry: tok.Expiry, AccessToken: tok.AccessToken, RefreshToken: tok.RefreshToken, TokenType: tok.TokenType}
	n, err := c.ctx.TokenAuthenticator.TokenSource(c.ctx.HTTPClient, token)

	return n, err
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func syncError(cid int64, tcid int64, r *http.Response, err error) {
	if r != nil {
		req, _ := httputil.DumpRequest(r.Request, true)
		res, _ := httputil.DumpResponse(r, true)
		e := models.SetTokenError(cid, tcid, r.StatusCode, r.Status, req, res)
		if e != nil {
			log.Println(e)
		}
	} else {
		e := models.SetTokenError(cid, tcid, 999, err.Error(), []byte{}, []byte{})
		if e != nil {
			log.Println(e)
		}
	}
	log.Printf("Contact Sync: %d %d %s", cid, tcid, err.Error())
}

func syncSuccess(cid int64, tcid int64, code int, status string) {
	models.SetTokenError(cid, tcid, code, status, []byte{}, []byte{})
}
