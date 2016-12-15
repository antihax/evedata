package eveConsumer

import (
	"encoding/json"
	"evedata/models"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	knownKills map[int64]bool
	mapLock    sync.RWMutex
	limiter    chan bool
)

func (c *EVEConsumer) initKillConsumer() {
	knownKills = make(map[int64]bool)
	k, err := models.GetKnownKillmails()
	if err != nil {
		log.Panic("Could not get known mails ", err)
	}
	for _, m := range k {
		knownKills[m] = true
	}
	limiter = make(chan bool, 10)
}

func (c *EVEConsumer) addKillmail(href string) error {

	// Check the kill is not known and early out if it is.
	hash := strings.Split(href, "/")[5]
	id, err := strconv.ParseInt(strings.Split(href, "/")[4], 10, 64)
	if err != nil {
		return err
	}
	mapLock.RLock()
	known := knownKills[id]
	mapLock.RUnlock()
	if known == true {
		return nil
	}

	limiter <- true
	go func(l chan bool, h string) error {
		defer func(l chan bool) { <-l }(l)

		kill, err := c.ctx.EVE.KillmailV1(h)
		if err != nil {
			return err
		}
		c.updateEntity(kill.Victim.Character.Href, kill.Victim.Character.ID)
		c.updateEntity(kill.Victim.Corporation.Href, kill.Victim.Corporation.ID)
		if kill.Victim.Alliance.ID != 0 {
			c.updateEntity(kill.Victim.Alliance.Href, kill.Victim.Alliance.ID)
		}
		models.AddKillmail(kill.KillID, kill.SolarSystem.ID, kill.KillTime.UTC(), kill.Victim.Character.ID,
			kill.Victim.Corporation.ID, kill.Victim.Alliance.ID, hash, kill.AttackerCount, kill.Victim.DamageTaken,
			kill.Victim.Position.X, kill.Victim.Position.Y, kill.Victim.Position.Z, kill.Victim.ShipType.ID,
			kill.War.ID)

		for _, item := range kill.Victim.Items {
			models.AddKillmailItems(kill.KillID, item.ItemType.ID, item.Flag, item.QuantityDestroyed,
				item.QuantityDropped, item.Singleton)
		}

		for _, attacker := range kill.Attackers {
			c.updateEntity(attacker.Character.Href, attacker.Character.ID)
			c.updateEntity(attacker.Corporation.Href, attacker.Corporation.ID)
			if attacker.Alliance.ID != 0 {
				c.updateEntity(attacker.Alliance.Href, attacker.Alliance.ID)
			}
			models.AddKillmailAttacker(kill.KillID, attacker.Character.ID, attacker.Corporation.ID, attacker.Alliance.ID,
				attacker.ShipType.ID, attacker.FinalBlow, attacker.DamageDone, attacker.WeaponType.ID,
				attacker.SecurityStatus)
		}
		mapLock.Lock()
		knownKills[id] = true
		mapLock.Unlock()
		return nil
	}(limiter, href)
	return nil
}

func (c *EVEConsumer) goZKillConsumer() error {
	type kill struct {
		Package struct {
			KillID int
			ZKB    struct {
				Hash string
			}
		}
	}

	for {
		k := kill{}

		err := c.getJSON("https://redisq.zkillboard.com/listen.php", &k)
		if err != nil {
			continue
		}
		if k.Package.KillID > 0 {
			c.addKillmail(fmt.Sprintf(c.ctx.EVE.GetCRESTURI()+"killmails/%d/%s/", k.Package.KillID, k.Package.ZKB.Hash))
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

	// Pull one date per minute.
	rate := time.Second * 60
	throttle := time.Tick(rate)

	for {
		<-throttle
		k := make(map[string]interface{})

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

		err := c.getJSON(fmt.Sprintf("https://zkillboard.com/api/history/%s/", date), &k)
		if err != nil {
			continue
		}

		for id, hash := range k {
			c.addKillmail(fmt.Sprintf(c.ctx.EVE.GetCRESTURI()+"killmails/%s/%s/", id, hash))
		}

		err = models.SetServiceState("zkilltemp", nextCheck, 1)
		if err != nil {
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
