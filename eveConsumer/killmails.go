package eveConsumer

import (
	"encoding/json"
	"evedata/models"
	"fmt"
	"log"
	"net/http"
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
	limiter <- true

	go func(l chan bool) error {
		defer func(l chan bool) { <-l }(l)

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

		kill, err := c.ctx.EVE.Killmail(href)
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
	}(limiter)
	return nil
}

func (c *EVEConsumer) goZKillConsumer() error {
	type kill struct {
		Package struct {
			KillID int

			ZKB struct {
				Hash string
			}
		}
	}

	for {
		k := kill{}

		err := getJSON("https://redisq.zkillboard.com/listen.php", &k)
		if err != nil {
			continue
		}
		if k.Package.KillID > 0 {
			c.addKillmail(fmt.Sprintf(c.ctx.EVE.GetCRESTURI()+"killmails/%d/%s/", k.Package.KillID, k.Package.ZKB.Hash))
		}
	}
}

func (c *EVEConsumer) goZKillTemporaryConsumer() error {
	r := struct {
		Value int
		Date  time.Time
	}{0, time.Now()}

	if err := c.ctx.Db.Get(&r, `
		SELECT value, date(nextCheck) AS date
			FROM states 
			WHERE state = 'zkilltemp'
			LIMIT 1;
		`); err != nil {
		return err
	}

	for {
		k := make(map[string]interface{})

		date := r.Date.Format("20060102")
		r.Date = r.Date.Add(time.Hour * 24)

		if r.Date.Sub(time.Now().UTC()) > 0 {
			r.Date = time.Now().UTC().Add(time.Hour * 24 * -365)
			log.Printf("Restart zKill Consumer to %s", r.Date.String())
		}

		err := getJSON(fmt.Sprintf("https://zkillboard.com/api/history/%s/", date), &k)
		if err != nil {
			continue
		}

		for id, hash := range k {
			c.addKillmail(fmt.Sprintf(c.ctx.EVE.GetCRESTURI()+"killmails/%s/%s/", id, hash))
		}

		_, err = c.ctx.Db.Exec("UPDATE states SET value = ?, nextCheck =? WHERE state = 'zkilltemp' LIMIT 1", 1, r.Date)
		if err != nil {
			continue
		}
	}
}

func getJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
