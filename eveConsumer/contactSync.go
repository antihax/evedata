package eveConsumer

import "log"

func (c *EveConsumer) contactSync() {
	rows, err := c.db.Query(
		`SELECT source, group_concat(destination)
			FROM contactSyncs GROUP BY source
		    HAVING max(nextSync) < UTC_TIMESTAMP()`)
	tx, err := c.db.Beginx()
	if err != nil {
		log.Printf("EVEConsumer: Failed starting transaction: %v", err)
		return
	}

	for rows.Next() {
		var (
			source int
			dest   string
		)

		err = rows.Scan(&source, &dest)
		//destinations := strings.Split(dest, ",")
		if err != nil {
			log.Printf("EVEConsumer: Failed Scanning Rows: %v", err)
			continue
		}
		char, err := c.eve.GetCharacterInfo(source)
		if err != nil {
			log.Printf("EVEConsumer: Failed getting character info %v", err)
			continue
		}

		var searchID int
		if char.AllianceID > 0 {
			searchID = char.AllianceID
		} else {
			searchID = char.CharacterID
		}

		// Active Wars
		// Would throw this into a procedure.. but cant use them with Golang sql...
		activeWars, err := c.db.Query(`
			SELECT defenderID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?;
		`, searchID, searchID, searchID, searchID)
		if err != nil {
			log.Printf("EVEConsumer: Failed Querying Active Wars: %v", err)
			continue
		}
		for activeWars.Next() {
			var id int

			err = activeWars.Scan(&id)
			if err != nil {
				log.Printf("EVEConsumer: Failed Scanning Active Wars: %v", err)
				continue
			}

		}
	}
	err = tx.Commit()
}
