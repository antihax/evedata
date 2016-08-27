package eveConsumer

import (
	"evedata/models"
	"log"
	"strings"
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
		`SELECT allianceid FROM alliance 
			WHERE cacheUntil < UTC_TIMESTAMP()`)
	if err != nil {
		return err
	}

	for alliances.Next() {
		var id int64
		err = alliances.Scan(&id)
		if err != nil {
			return err
		}
		a, err := c.ctx.EVE.AllianceByID(id)
		if err != nil {
			return err
		}
		err = models.UpdateAlliance(a.ID, a.Name, a.CorporationsCount, a.ShortName, a.ExecutorCorporation.ID,
			a.StartDate.UTC(), a.Deleted, a.Description, a.CreatorCorporation.ID, a.CreatorCharacter.ID, a.CacheUntil.UTC())
		if err != nil {
			return err
		}
		err = c.updateCharacter(a.CreatorCharacter.ID)
		if err != nil {
			return err
		}
	}
	alliances.Close()

	corporations, err := c.ctx.Db.Query(
		`SELECT corporationid FROM corporation 
			WHERE cacheUntil < UTC_TIMESTAMP()`)
	if err != nil {
		return err
	}

	for corporations.Next() {
		var id int64
		err = corporations.Scan(&id)
		if err != nil {
			return err
		}
		err = c.updateCorporation(id)
		if err != nil {
			return err
		}
	}
	corporations.Close()

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

	w, err := c.ctx.EVE.Alliances(0)

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
			err := c.updateAlliance(r.HRef)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EVEConsumer) updateEntity(href string, id int64) error {
	var err error
	if strings.Contains(href, "alliances") {
		err = c.updateAlliance(href)
	} else if strings.Contains(href, "corporations") {
		err = c.updateCorporation(id)
	} else if strings.Contains(href, "characters") {
		err = c.updateCharacter(id)
	}
	return err
}

func (c *EVEConsumer) updateAlliance(href string) error {
	a, err := c.ctx.EVE.Alliance(href)
	if err != nil {
		return err
	}
	err = models.AddCRESTRef(a.ID, href)
	if err != nil {
		return err
	}
	err = models.UpdateAlliance(a.ID, a.Name, a.CorporationsCount, a.ShortName, a.ExecutorCorporation.ID,
		a.StartDate.UTC(), a.Deleted, a.Description, a.CreatorCorporation.ID, a.CreatorCharacter.ID, a.CacheUntil.UTC())
	if err != nil {
		return err
	}
	err = c.updateCharacter(a.CreatorCharacter.ID)
	if err != nil {
		return err
	}

	for _, corp := range a.Corporations {
		err = c.updateCorporation(corp.ID)
		if err != nil {
			return err
		}
		err = models.AddCRESTRef(corp.ID, corp.Href)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *EVEConsumer) updateCorporation(id int64) error {
	a, err := c.ctx.EVE.GetCorporationPublicSheet(id)
	if err != nil {
		return err
	}

	err = models.UpdateCorporation(a.CorporationID, a.CorporationName, a.Ticker, a.CEOID, a.StationID,
		a.Description, a.AllianceID, a.FactionID, a.URL, a.MemberCount, a.Shares, a.CachedUntil.UTC())
	if err != nil {
		return err
	}

	err = c.updateCharacter(a.CEOID)
	if err != nil {
		return err
	}

	return nil
}

func (c *EVEConsumer) updateCharacter(id int64) error {
	a, err := c.ctx.EVE.GetCharacterInfo(id)
	if err != nil {
		return err
	}
	err = models.UpdateCharacter(a.CharacterID, a.CharacterName, a.BloodlineID, a.AncestryID, a.CorporationID, a.AllianceID, a.Race, a.SecurityStatus, a.CachedUntil.UTC())
	if err != nil {
		return err
	}
	return nil
}
