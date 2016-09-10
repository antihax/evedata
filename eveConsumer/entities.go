package eveConsumer

import (
	"evedata/models"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

func (c *EVEConsumer) checkAlliances() {

	err := c.collectEntitiesFromCREST()
	if err != nil {
		log.Printf("EVEConsumer: collecting entities: %v", err)
	}
	err = c.updateEntities()
	if err != nil {
		log.Printf("EVEConsumer: updating entities: %v", err)
	}

}

func (c *EVEConsumer) updateEntities() error {

	alliances, err := c.ctx.Db.Query(
		`SELECT allianceid, crestRef, cacheUntil FROM alliance A
INNER JOIN crestID C ON A.allianceID = C.id
			WHERE cacheUntil < UTC_TIMESTAMP()  
UNION
SELECT corporationid, crestRef, cacheUntil FROM corporation A
INNER JOIN crestID C ON A.corporationID = C.id
			WHERE cacheUntil < UTC_TIMESTAMP()
UNION
(SELECT characterID, crestRef, cacheUntil FROM eve.character A
INNER JOIN crestID C ON A.characterID = C.id
			WHERE cacheUntil < UTC_TIMESTAMP())
            
            ORDER BY cacheUntil ASC`)
	if err != nil {
		return err
	}

	for alliances.Next() {
		var (
			id      int64
			href    string
			nothing string
		)

		err = alliances.Scan(&id, &href, &nothing)
		if err != nil {
			return err
		}

		if err = c.updateEntity(href, id); err != nil {
			return err
		}

	}
	alliances.Close()

	return nil
}

func (c *EVEConsumer) collectEntitiesFromCREST() error {
	r := struct {
		Value int
		Wait  int
	}{0, 0}

	if err := c.ctx.Db.Get(&r, `
		SELECT value, TIME_TO_SEC(TIMEDIFF(nextCheck, UTC_TIMESTAMP())) AS wait
			FROM states 
			WHERE state = 'alliances'
			LIMIT 1;
		`); err != nil {
		return err
	}

	if r.Wait >= 0 {
		return nil
	}

	w, err := c.ctx.EVE.Alliances(r.Value)

	if err != nil {
		return err
	}

	// Loop through all of the pages
	for ; w != nil; w, err = w.NextPage() {
		if err != nil {
			return err
		}

		// Update state so we dont have two polling at once.
		_, err = c.ctx.Db.Exec("UPDATE states SET value = ?, nextCheck =? WHERE state = 'alliances' LIMIT 1", w.Page, w.CacheUntil)

		if err != nil {
			return err
		}

		for _, r := range w.Items {
			if err = c.updateEntity(r.HRef, r.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EVEConsumer) updateEntity(href string, id int64) error {
	var (
		err error
		t   time.Duration
	)
	go func() {
		r := c.ctx.Cache.Get()
		defer r.Close()
		i, err := redis.Bool(r.Do("EXISTS", "EVEDATA_entity:"+href))
		if err == nil || i == true {

			if strings.Contains(href, "alliances") {
				t, err = c.updateAlliance(href)
			} else if strings.Contains(href, "corporations") {
				t, err = c.updateCorporation(id)
			} else if strings.Contains(href, "characters") {
				t, err = c.updateCharacter(id)
			}
			if err != nil {
				return
			}
			err = models.AddCRESTRef(id, href)
			if err != nil {
				return
			}
			r.Do("SETEX", "EVEDATA_entity:"+href, (int64)(t.Seconds())+1, true)
		} else {
			return
		}

		return

	}()
	return err
}

func (c *EVEConsumer) updateAlliance(href string) (time.Duration, error) {

	a, err := c.ctx.EVE.Alliance(href)
	if err != nil {
		return 1, err
	}

	err = models.UpdateAlliance(a.ID, a.Name, a.CorporationsCount, a.ShortName, a.ExecutorCorporation.ID,
		a.StartDate.UTC(), a.Deleted, a.Description, a.CreatorCorporation.ID, a.CreatorCharacter.ID, a.CacheUntil.UTC())
	if err != nil {
		return 1, err
	}
	err = c.updateEntity(a.CreatorCharacter.Href, a.CreatorCharacter.ID)
	if err != nil {
		return 1, err
	}

	for _, corp := range a.Corporations {
		err = c.updateEntity(corp.Href, corp.ID)
		if err != nil {
			return 1, err
		}
	}
	t := a.CacheUntil.Sub(time.Now().UTC())
	return t, nil
}

func (c *EVEConsumer) updateCorporation(id int64) (time.Duration, error) {

	a, err := c.ctx.EVE.GetCorporationPublicSheet(id)
	if err != nil {
		return 1, err
	}

	err = models.UpdateCorporation(a.CorporationID, a.CorporationName, a.Ticker, a.CEOID, a.StationID,
		a.Description, a.AllianceID, a.FactionID, a.URL, a.MemberCount, a.Shares, a.CachedUntil.UTC())
	if err != nil {
		return 1, err
	}

	chref := "https://crest-tq.eveonline.com/" + fmt.Sprintf("characters/%d/", a.CEOID)
	err = c.updateEntity(chref, a.CEOID)
	if err != nil {
		return 1, err
	}

	t := a.CachedUntil.Sub(time.Now().UTC())
	return t, nil
}

func (c *EVEConsumer) updateCharacter(id int64) (time.Duration, error) {
	a, err := c.ctx.EVE.GetCharacterInfo(id)
	if err != nil {
		return 1, err
	}
	err = models.UpdateCharacter(a.CharacterID, a.CharacterName, a.BloodlineID, a.AncestryID, a.CorporationID, a.AllianceID, a.Race, a.SecurityStatus, a.CachedUntil.UTC())
	if err != nil {
		return 1, err
	}
	t := a.CachedUntil.Sub(time.Now().UTC())
	return t, nil
}
