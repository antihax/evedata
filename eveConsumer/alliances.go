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

	w, err := c.ctx.EVE.Alliances(r.Value)

	if err != nil {
		log.Printf("EVEConsumer: Error checking wars: %v", err)
		return
	}

	// Loop through all of the pages
	for ; w != nil; w, err = w.NextPage() {
		if err != nil {
			log.Printf("EVEConsumer: Could not start transaction for wars: %v", err)
			continue
		}

		// Update state so we dont have two polling at once.
		_, err = c.ctx.Db.Exec("UPDATE states SET value = ?, nextCheck =? WHERE state = 'alliances' LIMIT 1", w.Page, w.CacheUntil)

		if err != nil {
			log.Printf("EVEConsumer: Could not update war state: %v", err)
			continue
		}

		for _, r := range w.Items {
			if err != nil {
				log.Printf("EVEConsumer: Failed writing wars: %v", err)
				continue
			}

			err = models.AddCRESTRef(r.ID, r.Href)
			if err != nil {
				log.Printf("EVEConsumer: Failed writing CREST ref: %v", err)
				continue
			}
		}

		if err != nil {
			log.Printf("EVEConsumer: Failed writing war allies: %v", err)
			return
		}
	}
}
