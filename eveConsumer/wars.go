package eveConsumer

import (
	"evedata/models"
	"fmt"
	"log"
	"time"
)

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
		c.updateWar(fmt.Sprintf(c.ctx.EVE.GetCRESTURI()+"wars/%d/", id))
	}
	return nil
}

func (c *EVEConsumer) collectWarsFromCREST() error {
	nextCheck, page, err := models.GetServiceState("wars")
	if err != nil {
		return err
	} else if nextCheck.Before(time.Now()) {
		return nil
	}

	w, err := c.ctx.EVE.WarsV1((int)(page))
	if err != nil {
		return err
	}

	// Loop through all of the pages
	for ; w != nil; w, err = w.NextPage() {
		// Update state so we dont have two polling at once.
		err = models.SetServiceState("wars", w.CacheUntil, (int32)(w.Page))
		if err != nil {
			continue
		}

		for _, r := range w.Items {
			c.updateWar(r.HRef)
		}
	}
	return nil
}

func (c *EVEConsumer) updateWar(href string) error {
	war, err := c.ctx.EVE.WarV1(href)
	if err != nil {
		return err
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
		return err
	}

	err = c.updateEntity(war.Aggressor.HRef, war.Aggressor.ID)
	if err != nil {
		return err
	}

	err = c.updateEntity(war.Defender.HRef, war.Defender.ID)
	if err != nil {
		return err
	}

	for _, a := range war.Allies {
		_, err = c.ctx.Db.Exec(`INSERT IGNORE INTO warAllies (id, allyID) VALUES(?,?);`, war.ID, a.ID)
		if err != nil {
			return err
		}

		if err = c.updateEntity(a.HRef, a.ID); err != nil {
			return err
		}
	}

	kills, err := war.KillmailsV1()
	if err != nil {
		return err
	}
	for _, kills := range kills.Items {
		err := c.addKillmail(kills.HRef)
		if err != nil {
			return err
		}
	}

	return nil
}
