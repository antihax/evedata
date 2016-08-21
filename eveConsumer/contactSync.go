package eveConsumer

import (
	"evedata/eveapi"
	"evedata/models"
	"log"
	"strconv"
	"strings"
)

func (c *EVEConsumer) contactSync() {

	// Gather characters for update
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
			source int64
			dest   string
		)

		err = rows.Scan(&source, &dest)
		if err != nil {
			log.Printf("EVEConsumer: Failed scan: %v", err)
			continue
		}
		destinations := strings.Split(dest, ",")
		if err != nil {
			log.Printf("EVEConsumer: Failed Scanning Rows: %v", err)
			continue
		}

		char, err := c.ctx.EVE.GetCharacterInfo(source)
		if err != nil {
			log.Printf("EVEConsumer: Failed getting character info %v", err)
			continue
		}

		// Authenticated Clients
		clients := make(map[int64]*eveapi.AuthenticatedClient)

		// Get authenticated clients for our destinations
		for _, cidS := range destinations {
			cid, _ := strconv.ParseInt(cidS, 10, 64)
			a, err := c.getClient(source, cid)
			if err != nil {
				log.Printf("EVEConsumer: Failed getClient %v", err)
				continue
			}
			clients[cid] = a
			if err != nil {
				log.Printf("EVEConsumer: Failed Getting Contacts: %v", err)
				continue
			}
		}

		// Find the ID to search for wars.
		var searchID int
		if char.AllianceID > 0 {
			searchID = char.AllianceID
		} else {
			searchID = char.CharacterID
		}

		// Active Wars
		activeWars, err := c.ctx.Db.Query(`
			SELECT K.id, crestRef FROM
			(SELECT defenderID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN crestID C ON C.id = K.id
		`, searchID, searchID, searchID, searchID)
		if err != nil {
			log.Printf("EVEConsumer: Failed Getting Active Wars: %v", err)
			continue
		}
		defer activeWars.Close()

		// Pending Wars
		pendingWars, err := c.ctx.Db.Query(`
			SELECT K.id, crestRef FROM
			(SELECT defenderID AS id FROM wars WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM wars WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN crestID C ON C.id = K.id
		`, searchID, searchID, searchID, searchID)

		if err != nil {
			log.Printf("EVEConsumer: Failed Getting Pending Wars: %v", err)
			continue
		}
		defer pendingWars.Close()

		type toAdd struct {
			id       int64
			ref      string
			standing float64
		}
		contactsToAdd := make(map[int64]*toAdd)

		for activeWars.Next() {
			con := &toAdd{standing: -10}

			if err = activeWars.Scan(&con.id, &con.ref); err != nil {
				log.Printf("EVEConsumer: Failed Scanning Active Wars: %v", err)
				continue
			}
			contactsToAdd[con.id] = con
		}

		for pendingWars.Next() {
			con := &toAdd{standing: -5}

			if err = pendingWars.Scan(&con.id, &con.ref); err != nil {
				log.Printf("EVEConsumer: Failed Scanning Active Wars: %v", err)
				continue
			}
			contactsToAdd[con.id] = con
		}

		for _, client := range clients {

			if client == nil {
				continue
			}

			// Copy the contactsToAdd map
			toProcess := make(map[int64]*toAdd)
			for k, v := range contactsToAdd {
				toProcess[k] = v
			}

			// Get the clients current contacts
			con, err := client.GetContacts()
			if err != nil {
				log.Printf("EVEConsumer: Failed Getting Client Contacts: %v", err)
				continue
			}

			contactSync := &models.ContactSync{Source: source, Destination: client.GetCharacterID()}
			contactSync.Updated(con.CacheUntil)
			// Loop through all contact pages
			for ; con != nil; con, err = con.NextPage() {
				for _, contact := range con.Items {
					// skip anything > -0.4
					if contact.Standing > -0.4 {
						continue
					}

					add := toProcess[contact.Contact.ID]
					if add != nil {
						// Contact is already listed.
						if contact.Standing != add.standing {
							err = client.SetContact(add.id, add.ref, add.standing)
							if err != nil {
								log.Printf("EVEConsumer: Failed SetContact: %v", err)
								continue
							}
						}
						// Don't need to do anything to this contact.
						delete(toProcess, contact.Contact.ID)
					} else {
						// No longer at war... delete the contact
						err = client.DeleteContact(contact.Contact.ID, contact.Contact.Href)
						if err != nil {
							log.Printf("EVEConsumer: Failed DeleteContact: %v", err)
							continue
						}
					}
				}

				// Add the remaining contacts
				for _, contact := range toProcess {
					err = client.SetContact(contact.id, contact.ref, contact.standing)
					if err != nil {
						log.Printf("EVEConsumer: Failed SetContact: %v", err)
						continue
					}
				}
			}
		}
	}
}

func (c *EVEConsumer) getClient(characterID int64, tokenCharacterID int64) (*eveapi.AuthenticatedClient, error) {
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
	n := c.ctx.TokenAuthenticator.GetClientFromToken(c.ctx.HTTPClient, token)

	return n, nil
}
