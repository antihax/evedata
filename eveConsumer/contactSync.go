package eveConsumer

import (
	"evedata/eveapi"
	"evedata/models"
	"log"
	"strconv"
	"strings"
)

func (c *EveConsumer) contactSync() {
	return
	// Gather characters for update
	rows, err := c.ctx.Db.Query(
		`SELECT source, group_concat(destination)
			FROM contactSyncs GROUP BY source
		    HAVING max(nextSync) < UTC_TIMESTAMP()`)
	tx, err := c.ctx.Db.Beginx()
	if err != nil {
		log.Printf("EVEConsumer: Failed starting transaction: %v", err)
		return
	}

	// Loop updatable characters
	for rows.Next() {
		var (
			source int64
			dest   string
		)

		err = rows.Scan(&source, &dest)
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

		clients := make(map[int64]*eveapi.AuthenticatedClient)

		for _, cidS := range destinations {
			cid, _ := strconv.ParseInt(cidS, 10, 64)
			a, err := c.getClient(source, cid)
			if err != nil {
				log.Printf("EVEConsumer: Failed client %v", err)
				continue
			}
			clients[cid] = a
		}

		var searchID int
		if char.AllianceID > 0 {
			searchID = char.AllianceID
		} else {
			searchID = char.CharacterID
		}

		// Active Wars
		// Would throw this into a procedure.. but cant use them with Golang sql...
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
			log.Printf("EVEConsumer: Failed Querying Active Wars: %v", err)
			continue
		}

		for activeWars.Next() {
			var id int64

			err = activeWars.Scan(&id)
			if err != nil {
				log.Printf("EVEConsumer: Failed Scanning Active Wars: %v", err)
				continue
			}
			for _, client := range clients {
				client.SetContact(id, -10)
			}
			log.Printf("%d\n", id)
		}
	}
	err = tx.Commit()
}

func (c *EveConsumer) getClient(characterID int64, tokenCharacterID int64) (*eveapi.AuthenticatedClient, error) {
	tok := models.CRESTToken{}

	if err := c.ctx.Db.QueryRowx(
		`SELECT  expiry, tokenType , accessToken , refreshToken, tokenCharacterID, characterID
			FROM crestTokens
			WHERE characterID = ? and tokenCharacterID = ?
			LIMIT 1`,
		characterID, tokenCharacterID).StructScan(&tok); err != nil {

		return nil, err
	}

	token := &eveapi.CRESTToken{Expiry: tok.Expiry, AccessToken: tok.AccessToken, RefreshToken: tok.RefreshToken, TokenType: tok.TokenType}
	n := c.ctx.TokenAuthenticator.GetClientFromToken(c.ctx.HTTPClient, token)

	return n, nil
}
