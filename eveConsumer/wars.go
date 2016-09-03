package eveConsumer

import "log"

func (c *EVEConsumer) checkWars() {
	err := c.collectWarsFromCREST()
	if err != nil {
		log.Printf("EVEConsumer: collecting wars: %v", err)
	}

	err = c.updateWars()
	if err != nil {
		log.Printf("EVEConsumer: updating wars: %v", err)
	}
}

func (c *EVEConsumer) updateWars() error {
	rows, err := c.ctx.Db.Query(
		`SELECT id FROM eve.wars 
			WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL) 
			AND cacheUntil < UTC_TIMESTAMP()`)
	if err != nil {

		return err
	}
	defer rows.Close()

	for rows.Next() {

		var id int
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		war, err := c.ctx.EVE.WarByID(id)
		if err != nil {
			return err
		}

		_, err = c.ctx.Db.Exec(`UPDATE wars SET
					timeFinished=?, 
					openForAllies=?, 
					mutual=?, 
					cacheUntil=? WHERE id = ?;`,
			war.TimeFinished.String(),
			war.OpenForAllies,
			war.Mutual,
			war.CacheUntil, war.ID)
		if err != nil {
			return err
		}

		for _, a := range war.Allies {
			_, err = c.ctx.Db.Exec(`INSERT IGNORE INTO warAllies (id, allyID) VALUES(?,?);`, war.ID, a.ID)
			if err != nil {
				return err
			}

			c.updateEntity(a.HRef, a.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EVEConsumer) collectWarsFromCREST() error {
	r := struct {
		Value int
		Wait  int
	}{0, 0}

	if err := c.ctx.Db.Get(&r, `
		SELECT value, TIME_TO_SEC(TIMEDIFF(nextCheck, UTC_TIMESTAMP())) AS wait
			FROM states 
			WHERE state = 'wars'
			LIMIT 1;
		`); err != nil {
		return err
	}

	if r.Wait >= 0 {
		//return nil
	}

	w, err := c.ctx.EVE.Wars(r.Value)

	if err != nil {
		return err
	}

	// Loop through all of the pages
	for ; w != nil; w, err = w.NextPage() {

		// Update state so we dont have two polling at once.
		_, err = c.ctx.Db.Exec("UPDATE states SET value = ?, nextCheck =? WHERE state = 'wars' LIMIT 1", w.Page, w.CacheUntil)

		if err != nil {
			continue
		}

		for _, r := range w.Items {

			war, err := c.ctx.EVE.War(r.HRef)
			if err != nil {
				continue
			}

			_, err = c.ctx.Db.Exec(`INSERT INTO wars
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
				continue
			}

			err = c.updateEntity(war.Aggressor.HRef, war.Aggressor.ID)
			if err != nil {
				continue
			}

			err = c.updateEntity(war.Defender.HRef, war.Defender.ID)
			if err != nil {
				continue
			}

			for _, a := range war.Allies {
				_, err = c.ctx.Db.Exec(`INSERT IGNORE INTO warAllies (id, allyID) VALUES(?,?);`, war.ID, a.ID)
				if err != nil {
					continue
				}

				if err = c.updateEntity(a.HRef, a.ID); err != nil {
					continue
				}
			}

			kills, err := war.GetKillmails()
			if err != nil {
				return err
			}
			for _, kills := range kills.Items {
				err := c.addKillmail(kills.HRef)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
