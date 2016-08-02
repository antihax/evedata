package eveConsumer

import "log"

func (c *EveConsumer) checkWars() {
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
				_, err = tx.Exec(`INSERT INTO warAllies
				(id, allyID) VALUES(?,?) ON DUPLICATE IGNORE;`,
					war.ID, a.ID)
				if err != nil {
					log.Printf("EVEConsumer: Failed writing war allies: %v", err)
					return
				}
			}
		}
		tx.Exec("UPDATE states SET value = ?, nextCheck =? WHERE state = 'wars' LIMIT 1", w.Page, w.CacheUntil)
		err = tx.Commit()
		if err != nil {
			log.Printf("EVEConsumer: Could not commit transaction for wars: %v", err)
			return
		}
	}

}
