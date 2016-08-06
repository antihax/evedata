package eveConsumer

import "log"

func (c *EveConsumer) checkWars() {
	c.collectWarsFromCREST()
	c.updateWars()
}

func (c *EveConsumer) updateWars() {
	rows, err := c.db.Query(
		`SELECT id FROM eve.wars 
			WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL) 
			AND cacheUntil < UTC_TIMESTAMP()`)
	tx, err := c.db.Beginx()
	if err != nil {
		log.Printf("EVEConsumer: Failed starting transaction: %v", err)
		return
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Printf("EVEConsumer: Failed translating war ID: %v", err)
			return
		}
		war, err := c.eve.WarByID(id)
		if err != nil {
			log.Printf("EVEConsumer: Failed reading wars: %v", err)
			return
		}

		_, err = tx.Exec(`UPDATE wars SET
					timeFinished=?, 
					openForAllies=?, 
					mutual=?, 
					cacheUntil=? WHERE id = ?;`,
			war.TimeFinished.String(),
			war.OpenForAllies,
			war.Mutual,
			war.CacheUntil, war.ID)

		if err != nil {
			log.Printf("EVEConsumer: Failed writing wars: %v", err)
			return
		}

		for _, a := range war.Allies {
			_, err = tx.Exec(`INSERT IGNORE INTO warAllies
				(id, allyID) VALUES(?,?);`,
				war.ID, a.ID)
			if err != nil {
				log.Printf("EVEConsumer: Failed writing war allies: %v", err)
				return
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("EVEConsumer: Failed writing war allies: %v", err)
		return
	}
}

func (c *EveConsumer) collectWarsFromCREST() {
	r := struct {
		Value int
		Wait  int
	}{0, 0}

	if err := c.db.Get(&r, `
		SELECT value, TIME_TO_SEC(TIMEDIFF(nextCheck, UTC_TIMESTAMP())) AS wait
			FROM states 
			WHERE state = 'wars'
			LIMIT 1;
		`); err != nil {
		log.Printf("EVEConsumer: Error checking state: %v", err)
		return
	}

	if r.Wait >= 0 {
		return
	}

	w, err := c.eve.Wars(r.Value)

	if err != nil {
		log.Printf("EVEConsumer: Error checking wars: %v", err)
		return
	}

	// Loop through all of the pages
	for ; w != nil; w, err = w.NextPage() {
		tx, err := c.db.Beginx()
		if err != nil {
			log.Printf("EVEConsumer: Could not start transaction for wars: %v", err)
			return
		}

		// Update state so we dont have two polling at once.
		_, err = c.db.Exec("UPDATE states SET value = ?, nextCheck =? WHERE state = 'wars' LIMIT 1", w.Page, w.CacheUntil)

		if err != nil {
			log.Printf("EVEConsumer: Could not update war state: %v", err)
			return
		}

		for _, r := range w.Items {
			war, err := c.eve.War(r.HRef)
			if err != nil {
				log.Printf("EVEConsumer: Failed reading wars: %v", err)
				return
			}

			_, err = tx.Exec(`INSERT INTO wars
				(id, timeFinished,timeStarted,timeDeclared,openForAllies,cacheUntil,aggressorID,defenderID,mutual)
				VALUES(?,?,?,?,?,?,?,?,?)
				ON DUPLICATE KEY UPDATE 
					timeFinished=VALUES(timeFinished), 
					openForAllies=VALUES(openForAllies), 
					mutual=VALUES(mutual), 
					cacheUntil=VALUES(cacheUntil);`,
				war.ID, war.TimeFinished.String(), war.TimeStarted.String(), war.TimeDeclared.String(),
				war.OpenForAllies, war.CacheUntil, war.Aggressor.ID,
				war.Defender.ID, war.Mutual)
			if err != nil {
				log.Printf("EVEConsumer: Failed writing wars: %v", err)
				return
			}
			for _, a := range war.Allies {
				_, err = tx.Exec(`INSERT IGNORE INTO warAllies
				(id, allyID) VALUES(?,?);`,
					war.ID, a.ID)
				if err != nil {
					log.Printf("EVEConsumer: Failed writing war allies: %v", err)
					return
				}
			}
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("EVEConsumer: Failed writing war allies: %v", err)
			return
		}
	}
}
