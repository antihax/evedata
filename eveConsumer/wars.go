package eveConsumer

import (
	"time"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("wars", warConsumer, "EVEDATA_warQueue")
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

	rows, err := c.ctx.Db.Query(
		`SELECT id FROM evedata.wars WHERE timeFinished = '0001-01-01 00:00:00'
			AND cacheUntil < UTC_TIMESTAMP();`)
	if err != nil {
		return err
	}
	defer rows.Close()

	r := c.ctx.Cache.Get()
	defer r.Close()

	for rows.Next() {
		var id int32
		err = rows.Scan(&id)
		if err != nil {
			r.Flush()
			return err
		}
		r.Send("SREM", "EVEDATA_knownFinishedWars", id)
		r.Send("SADD", "EVEDATA_warQueue", id)
	}

	return r.Flush()
}

/*
// CCP disabled ESI wars. Go back to CREST until fixed.
func (c *EVEConsumer) collectWarsFromCREST() error {
	nextCheck, page, err := models.GetServiceState("wars")
	if err != nil {
		return err
	} else if nextCheck.After(time.Now()) {
		return nil
	}

	w, err := c.ctx.ESI.EVEAPI.WarsV1((int)(page))
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
}*/

func (c *EVEConsumer) collectWarsFromCREST() error {

	// Get the update state
	nextCheck, _, err := models.GetServiceState("wars")
	if err != nil {
		return err
	} else if nextCheck.After(time.Now().UTC()) { // Check if the cache timer has expired
		return nil
	}
	// Loop through all pages
	wars, res, err := c.ctx.ESI.ESI.WarsApi.GetWars(nil, nil)
	if err != nil {
		return err
	}

	r := c.ctx.Cache.Get()
	defer r.Close()
	for _, id := range wars {
		r.Send("SADD", "EVEDATA_warQueue", id)
	}
	r.Flush()
	models.SetServiceState("wars", goesi.CacheExpires(res), 1)

	return nil
}

func warConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
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

	// Get the war information
	war, res, err := c.ctx.ESI.ESI.WarsApi.GetWarsWarId(nil, (int32)(v), nil)
	if err != nil {
		r.Do("SADD", "EVEDATA_warQueue", v)
		return false, err
	}

	// save the aggressor id
	var aggressor, defender int32
	if war.Aggressor.AllianceId > 0 {
		aggressor = war.Aggressor.AllianceId
		EntityAllianceAddToQueue(aggressor, &r)
	} else {
		aggressor = war.Aggressor.CorporationId
		EntityCorporationAddToQueue(aggressor, &r)
	}

	// save the defender id
	if war.Defender.AllianceId > 0 {
		defender = war.Defender.AllianceId
		EntityAllianceAddToQueue(defender, &r)
	} else {
		defender = war.Defender.CorporationId
		EntityCorporationAddToQueue(defender, &r)
	}

	_, err = models.RetryExec(`INSERT INTO evedata.wars
				(id, timeFinished,timeStarted,timeDeclared,openForAllies,cacheUntil,aggressorID,defenderID,mutual)
				VALUES(?,?,?,?,?,?,?,?,?)
				ON DUPLICATE KEY UPDATE 
					timeFinished=VALUES(timeFinished), 
					openForAllies=VALUES(openForAllies), 
					mutual=VALUES(mutual), 
					cacheUntil=VALUES(cacheUntil);`,
		war.Id, war.Finished.Format(models.SQLTimeFormat), war.Started, war.Declared,
		war.OpenForAllies, goesi.CacheExpires(res), aggressor,
		defender, war.Mutual)
	if err != nil {
		r.Do("SADD", "EVEDATA_warQueue", v)
		return false, err
	}

	// Add information on allies in the war
	for _, a := range war.Allies {
		var ally int32
		if a.AllianceId > 0 {
			ally = a.AllianceId
			EntityAllianceAddToQueue(ally, &r)
		} else {
			ally = a.CorporationId
			EntityCorporationAddToQueue(ally, &r)
		}

		_, err = c.ctx.Db.Exec(`INSERT INTO evedata.warAllies (id, allyID) VALUES(?,?) ON DUPLICATE KEY UPDATE id = id;`, war.Id, ally)
		if err != nil {
			return false, err
		}
	}

	// If the war ended, cache the ID in redis to prevent needlessly pulling
	if war.Finished.IsZero() == false && war.Finished.Before(time.Now().UTC()) {
		r.Do("SADD", "EVEDATA_knownFinishedWars", war.Id)
	}

	// Loop through all the killmail pages
	for i := 1; ; i++ {
		kills, _, err := c.ctx.ESI.ESI.WarsApi.GetWarsWarIdKillmails(nil, war.Id, map[string]interface{}{"page": (int32)(i)})
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
