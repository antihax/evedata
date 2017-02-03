package eveConsumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("killmails", killmailsConsumer, "EVEDATA_killQueue")
}

func killmailsConsumer(c *EVEConsumer, r redis.Conn) (bool, error) {
	ret, err := r.Do("SPOP", "EVEDATA_killQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}

	v, err := redis.String(ret, err)
	if err != nil {
		return false, err
	}

	// split id:hash
	split := strings.Split(v, ":")
	if len(split) != 2 {
		return false, errors.New("string must be id:hash")
	}
	// convert ID to int64
	id, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return false, err
	}

	// We know this kill. Early out.
	i, err := redis.Int(r.Do("SISMEMBER", "EVEDATA_knownKills", (int32)(id)))
	if err == nil && i == 1 {
		return false, err
	}

	err = c.killmailGetAndSave((int32)(id), split[1])
	if err != nil {
		return false, err
	}
	return true, err

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
	r := c.ctx.Cache.Get()
	defer r.Close()
	r.Do("SADD", "EVEDATA_knownKills", id)
	return nil
}

// Launched in go routine
func (c *EVEConsumer) initKillConsumer(r redis.Conn) {
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
}

// Go get the killmail from CCP. Called from the queue consumer.
func (c *EVEConsumer) killmailGetAndSave(id int32, hash string) error {
	// Get the killmail from CCP
	kill, r, err := c.ctx.ESI.V1.KillmailsApi.GetKillmailsKillmailIdKillmailHash(id, hash, nil)

	// If we get a 500 error, add the mail back to the queue so we can try again later.
	if r != nil {
		if r.StatusCode >= 500 {
			return err
		}
	}

	if err != nil {
		return err
	}

	save := true
	old := time.Now().Add(time.Hour * -(24 * 365))
	if kill.KillmailTime.UTC().Before(old) {
		save = false
	}

	redis := c.ctx.Cache.Get()
	defer redis.Close()

	EntityAddToQueue(kill.Victim.CharacterId, redis)
	EntityAddToQueue(kill.Victim.CorporationId, redis)
	if kill.Victim.AllianceId != 0 {
		EntityAddToQueue(kill.Victim.AllianceId, redis)
	}
	if save {
		err = models.AddKillmail(kill.KillmailId, kill.SolarSystemId, kill.KillmailTime.UTC(), kill.Victim.CharacterId,
			kill.Victim.CorporationId, kill.Victim.AllianceId, hash, len(kill.Attackers), kill.Victim.DamageTaken,
			kill.Victim.Position.X, kill.Victim.Position.Y, kill.Victim.Position.Z, kill.Victim.ShipTypeId,
			kill.WarId)
		if err != nil {
			return err
		}

		for _, item := range kill.Victim.Items {
			err = models.AddKillmailItems(kill.KillmailId, item.ItemTypeId, item.Flag, item.QuantityDestroyed,
				item.QuantityDropped, item.Singleton)
			if err != nil {
				return err
			}
		}
	}
	for _, attacker := range kill.Attackers {
		EntityAddToQueue(attacker.CharacterId, redis)
		EntityAddToQueue(attacker.CorporationId, redis)
		if attacker.AllianceId != 0 {
			EntityAddToQueue(attacker.AllianceId, redis)
		}
		if save {
			err = models.AddKillmailAttacker(kill.KillmailId, attacker.CharacterId, attacker.CorporationId, attacker.AllianceId,
				attacker.ShipTypeId, attacker.FinalBlow, attacker.DamageDone, attacker.WeaponTypeId,
				attacker.SecurityStatus)
			if err != nil {
				return err
			}
		}
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

	time.Sleep(time.Second * 5)

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

	// three per second until we catchup
	rate := time.Second * 60 // ((60 * 60 * 24) / 365)
	throttle := time.Tick(rate)

	for {
		<-throttle

		// Move to the next day
		date := nextCheck.Format("20060102")
		nextCheck = nextCheck.Add(time.Hour * 24)

		// If we are at today, restart from 90 days
		if nextCheck.Sub(time.Now().UTC()) > 0 {
			nextCheck = time.Now().UTC().Add(time.Hour * 24 * -365)
			log.Printf("Delete old killmails")
			models.MaintKillMails()

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
