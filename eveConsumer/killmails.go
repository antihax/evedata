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
)

var (
	knownKills map[int64]bool
	mapLock    sync.RWMutex
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
}

func (c *EVEConsumer) addKillmail(href string) error {
	go func() error {
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
	}()
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

func getJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
