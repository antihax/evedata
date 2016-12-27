package eveConsumer

import (
	"encoding/json"
	"evedata/models"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Add a killmail to the queue
func (c *EVEConsumer) killmailGoQueueConsumer() error {
	r := c.ctx.Cache.Get()
	defer r.Close()
	for {
		// Pop off an element from the queue.
		ret, err := r.Do("SPOP", "EVEDATA_killQueue")
		if ret == nil {
			// Let's sleep for a bit.
			time.Sleep(time.Second * 20)
			continue
		}
		v, err := redis.String(ret, err)

		split := strings.Split(v, ":")
		if len(split) != 2 {
			log.Printf("Killmail error: invalid format\n", err)
			continue
		}
		id, err := strconv.ParseInt(split[0], 10, 32)
		if err != nil {
			log.Printf("Killmail error: %v\n", err)
			continue
		}
		// We know this kill. Early out.
		i, err := redis.Int(r.Do("SISMEMBER", "EVEDATA_knownKills", id))
		if err == nil && i == 1 {
			continue
		}
		err = c.killmailGetAndSave((int32)(id), split[1])
		if err != nil {
			log.Printf("Killmail error: %v\n", err)
			continue
		}

	}
}

// Add a killmail to the queue
func (c *EVEConsumer) killmailAddToQueue(id int32, hash string) error {
	r := c.ctx.Cache.Get()
	defer r.Close()
	key := fmt.Sprintf("%d:%s", id, hash)

	// We know this kill. Early out.
	i, err := redis.Int(r.Do("SISMEMBER", "EVEDATA_knownKills", id))
	if err == nil && i == 1 {
		return err
	}

	// Add the mail to the queue
	_, err = r.Do("SADD", "EVEDATA_killQueue", key)
	return err
}

// Say we know this killmail
func (c *EVEConsumer) killmailSetKnown(id int32) error {
	go func() {
		r := c.ctx.Cache.Get()
		defer r.Close()
		r.Do("SADD", "EVEDATA_knownKills", id)
	}()
	return nil
}

// Launched in go routine
func (c *EVEConsumer) initKillConsumer() {
	// Get a redis connection from the pool
	r := c.ctx.Cache.Get()
	defer r.Close()

	// get the list of know killmails
	k, err := models.GetKnownKillmails()
	if err != nil {
		log.Panic("Could not get known killmails ", err)
	}

	// Build a pipeline request to add the killmail IDs to redis
	for _, m := range k {
		r.Send("SADD", "EVEDATA_knownKills", m)
	}

	// Send the request to add
	r.Flush()

	log.Printf("Loaded %d known killmails\n", len(k))

	// Start the killmail queue consumer.
	for i := 0; i < 25; i++ {
		go c.killmailGoQueueConsumer()
	}
}

// Go get the killmail from CCP. Called from the queue consumer.
func (c *EVEConsumer) killmailGetAndSave(id int32, hash string) error {
	// Get the killmail from CCP
	kill, _, err := c.ctx.ESI.KillmailsApi.GetKillmailsKillmailIdKillmailHash(id, hash, nil)
	if err != nil {
		return err
	}
	c.updateESIEntitys(kill.Victim.CharacterId)
	c.updateESIEntitys(kill.Victim.CorporationId)
	if kill.Victim.AllianceId != 0 {
		c.updateESIEntitys(kill.Victim.AllianceId)
	}
	models.AddKillmail(kill.KillmailId, kill.SolarSystemId, kill.KillmailTime.UTC(), kill.Victim.CharacterId,
		kill.Victim.CorporationId, kill.Victim.AllianceId, hash, len(kill.Attackers), kill.Victim.DamageTaken,
		kill.Victim.Position.X, kill.Victim.Position.Y, kill.Victim.Position.Z, kill.Victim.ShipTypeId,
		kill.WarId)

	for _, item := range kill.Victim.Items {
		models.AddKillmailItems(kill.KillmailId, item.ItemTypeId, item.Flag, item.QuantityDestroyed,
			item.QuantityDropped, item.Singleton)
	}

	for _, attacker := range kill.Attackers {
		c.updateESIEntitys(attacker.CharacterId)
		c.updateESIEntitys(attacker.CorporationId)
		if attacker.AllianceId != 0 {
			c.updateESIEntitys(attacker.AllianceId)
		}
		models.AddKillmailAttacker(kill.KillmailId, attacker.CharacterId, attacker.CorporationId, attacker.AllianceId,
			attacker.ShipTypeId, attacker.FinalBlow, attacker.DamageDone, attacker.WeaponTypeId,
			attacker.SecurityStatus)
	}
	c.killmailSetKnown((int32)(id))
	return nil
}

// Collect killmails from RedisQ (zkillboard live feed)
func (c *EVEConsumer) goZKillConsumer() error {
	type kill struct {
		Package struct {
			KillID int32
			ZKB    struct {
				Hash string
			}
		}
	}

	for {
		k := kill{}
		err := c.getJSON("https://redisq.zkillboard.com/listen.php", &k)
		if err != nil {
			log.Printf("Zkill error: %v\n", err)
			continue
		}
		if k.Package.KillID > 0 {
			c.killmailAddToQueue(k.Package.KillID, k.Package.ZKB.Hash)
		}
	}
}

// Go Routine to collect killmails from ZKill API.
// Loops collecting one year of kill mails.
func (c *EVEConsumer) goZKillTemporaryConsumer() error {
	// Start from where we left off.
	nextCheck, _, err := models.GetServiceState("zkilltemp")
	if err != nil {
		return err
	}

	// Spread out over a day.
	rate := time.Second * ((60 * 60 * 24) / 365)
	throttle := time.Tick(rate)

	for {
		<-throttle

		// Move to the next day
		date := nextCheck.Format("20060102")
		nextCheck = nextCheck.Add(time.Hour * 24)

		// If we are at today, restart from one year ago
		if nextCheck.Sub(time.Now().UTC()) > 0 {
			nextCheck = time.Now().UTC().Add(time.Hour * 24 * -365)
			log.Printf("Delete old killmails")
			c.ctx.Db.Exec("CALL removeOldKillmails();")

			log.Printf("Restart zKill Consumer to %s", nextCheck.String())
		}

		// Get the kill history from ZKill for this day.
		k := make(map[string]interface{})
		err := c.getJSON(fmt.Sprintf("https://zkillboard.com/api/history/%s/", date), &k)
		if err != nil {
			continue
		}

		// Loop through the killmails
		for idS, hash := range k {
			id, err := strconv.ParseInt(idS, 10, 32)
			if err != nil {
				log.Printf("Zkill Consumer Error: %v", err)
				continue
			}

			// Add to the killmail queue
			err = c.killmailAddToQueue((int32)(id), hash.(string))
			if err != nil {
				log.Printf("Zkill Consumer Error: %v", err)
				continue
			}
		}

		// Update the current state on the database so we can restart where we left off
		err = models.SetServiceState("zkilltemp", nextCheck, 1)
		if err != nil {
			log.Printf("Zkill Consumer Error: %v", err)
			continue
		}
	}
}

func (c *EVEConsumer) getJSON(url string, target interface{}) error {
	r, err := c.ctx.HTTPClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
