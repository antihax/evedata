package eveConsumer

import (
	"evedata/models"
	"log"
)

func (c *EVEConsumer) checkAlliances() {
	c.collectAlliancesFromCREST()

}

func (c *EVEConsumer) collectAlliancesFromCREST() {
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
		log.Printf("EVEConsumer: Error checking state: %v", err)
		return
	}

	if r.Wait >= 0 {
		return
	}

	w, err := c.ctx.EVE.Alliances(0)

	if err != nil {
		log.Printf("EVEConsumer: Error checking Alliances: %v", err)
		return
	}

	// Loop through all of the pages
	for ; w != nil; w, err = w.NextPage() {
		if err != nil {
			log.Printf("EVEConsumer: Could not start transaction for Alliances: %v", err)
			continue
		}

		// Update state so we dont have two polling at once.
		_, err = c.ctx.Db.Exec("UPDATE states SET value = ?, nextCheck =? WHERE state = 'alliances' LIMIT 1", w.Page, w.CacheUntil)

		if err != nil {
			log.Printf("EVEConsumer: Could not update alliance state: %v", err)
			continue
		}

		for _, r := range w.Items {
			err := c.updateAlliance(r.HRef)
			if err != nil {
				log.Printf("EVEConsumer: Failed writing updating alliance: %v", err)
				continue
			}
		}

	}
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
		a.StartDate.UTC(), a.Deleted, a.Description, a.CreatorCorporation.ID, a.CreatorCharacter.ID)
	if err != nil {
		return err
	}

	return nil
}
