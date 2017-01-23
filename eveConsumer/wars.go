package eveConsumer

import (
	"fmt"
	"time"

	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("wars", warConsumer)
	addTrigger("wars", warsTrigger)
}

func warsTrigger(c *EVEConsumer) (bool, error) {

	err := c.collectWarsFromCREST()
	if err != nil {
		return false, err
	}

	err = c.updateWars()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *EVEConsumer) warAddToQueue(id int32) error {
	r := c.ctx.Cache.Get()
	defer r.Close()

	// This war is over. Early out.
	i, err := redis.Int(r.Do("SISMEMBER", "EVEDATA_knownFinishedWars", id))
	if err == nil && i == 1 {
		return err
	}

	// Add the war to the queue
	_, err = r.Do("SADD", "EVEDATA_warQueue", id)
	return err
}

func (c *EVEConsumer) updateWars() error {
	r := c.ctx.Cache.Get()
	defer r.Close()
	rows, err := c.ctx.Db.Query(
		`SELECT id FROM evedata.wars WHERE timeFinished = '0001-01-01 00:00:00'
			AND cacheUntil < UTC_TIMESTAMP();`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int32
		err = rows.Scan(&id)
		if err != nil {
			return err
		}

		r.Do("SREM", "EVEDATA_knownFinishedWars", id)
		c.warAddToQueue(id)
	}
	return nil
}

// CCP disabled ESI wars. Go back to CREST until fixed.
func (c *EVEConsumer) collectWarsFromCREST() error {
	nextCheck, page, err := models.GetServiceState("wars")
	if err != nil {
		return err
	} else if nextCheck.After(time.Now()) {
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
			c.warAddToQueue((int32)(r.ID))
		}
	}
	return nil
}

/*func (c *EVEConsumer) collectWarsFromCREST() error {
	nextCheck, _, err := models.GetServiceState("wars")
	if err != nil {
		return err
	} else if nextCheck.After(time.Now()) {
		return nil
	}

	var page int32 = 1

	for {
		wars, _, err := c.ctx.ESI.WarsApi.GetWars(map[string]interface{}{"page": page})
		if err != nil {
			return err
		} else if len(wars) == 0 { // end of the pages
			break
		}
		for _, id := range wars {
			c.warAddToQueue(id)
		}
		page++
	}

	models.SetServiceState("wars", time.Now().UTC().Add(time.Hour), 1)
	return nil
}*/

func warConsumer(c *EVEConsumer, r redis.Conn) (bool, error) {
	ret, err := r.Do("SPOP", "EVEDATA_warQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}

	v, err := redis.Int(ret, err)
	if err != nil {
		return false, err
	}
	href := fmt.Sprintf(c.ctx.EVE.GetCRESTURI()+"wars/%d/", v)
	war, err := c.ctx.EVE.WarV1(href)
	if err != nil {
		return false, err
	}
	_, err = models.RetryExec(`INSERT INTO evedata.wars
				(id, timeFinished,timeStarted,timeDeclared,openForAllies,cacheUntil,aggressorID,defenderID,mutual)
				VALUES(?,?,?,?,?,?,?,?,?)
				ON DUPLICATE KEY UPDATE 
					timeFinished=VALUES(timeFinished), 
					openForAllies=VALUES(openForAllies), 
					mutual=VALUES(mutual), 
					cacheUntil=VALUES(cacheUntil);`,
		war.ID, war.TimeFinished.Format(models.SQLTimeFormat), war.TimeStarted.Format(models.SQLTimeFormat), war.TimeDeclared.Format(models.SQLTimeFormat),
		war.OpenForAllies, war.CacheUntil, war.Aggressor.ID,
		war.Defender.ID, war.Mutual)
	if err != nil {
		return false, err
	}

	err = EntityAddToQueue((int32)(war.Aggressor.ID), r)
	if err != nil {
		return false, err
	}

	err = EntityAddToQueue((int32)(war.Defender.ID), r)
	if err != nil {
		return false, err
	}

	for _, a := range war.Allies {
		_, err = c.ctx.Db.Exec(`INSERT IGNORE INTO evedata.warAllies (id, allyID) VALUES(?,?);`, war.ID, a.ID)
		if err != nil {
			return false, err
		}

		if err = EntityAddToQueue((int32)(a.ID), r); err != nil {
			return false, err
		}
	}

	if war.TimeFinished.UTC().Before(time.Now()) {
		r.Do("SADD", "EVEDATA_knownFinishedWars", (int32)(war.ID))
	}

	// Loop through all the killmail pages
	for i := 1; ; i++ {
		kills, _, err := c.ctx.ESI.WarsApi.GetWarsWarIdKillmails((int32)(war.ID), map[string]interface{}{"page": (int32)(i)})
		if err != nil {
			return false, err
		}

		// No more kills to get, let`s get out of the loop.
		if len(kills) == 0 {
			break
		}

		// Add mails to the queue (queue will handle known mails)
		for _, k := range kills {
			err := c.killmailAddToQueue(k.KillmailId, k.KillmailHash)
			if err != nil {
				return false, err
			}
		}
	}

	// All good bro.
	return true, nil
}
