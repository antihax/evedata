package eveConsumer

import (
	"evedata/models"
	"log"
)

func (c *EVEConsumer) checkNPCCorps() {
	err := c.collectNPCCorps()
	if err != nil {
		log.Printf("EVEConsumer: collecting loyalty Point Store Items: %v", err)
	}
}

func (c *EVEConsumer) collectNPCCorps() error {
	r := struct {
		Value int
		Wait  int
	}{0, 0}

	if err := c.ctx.Db.Get(&r, `
		SELECT value, TIME_TO_SEC(TIMEDIFF(nextCheck, UTC_TIMESTAMP())) AS wait
			FROM states 
			WHERE state = 'npcCorps'
			LIMIT 1;
		`); err != nil {
		return err
	}

	if r.Wait >= 0 {
		return nil
	}

	w, err := c.ctx.EVE.NPCCorporations(1)
	if err != nil {
		return err
	}

	// Loop through all of the pages
	for ; w != nil; w, err = w.NextPage() {
		// Update state so we dont have two polling at once.
		_, err = c.ctx.Db.Exec("UPDATE states SET nextCheck = DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY) WHERE state = 'npcCorps' LIMIT 1")
		if err != nil {
			return err
		}

		for _, npcCorp := range w.Items {
			if npcCorp.LoyaltyStore.Href == "" {
				continue
			}

			c.updateEntity(npcCorp.Href, npcCorp.ID)
			store, err := c.ctx.EVE.LoyaltyPointStore(npcCorp.LoyaltyStore.Href)
			if err != nil {
				continue
			}

			for ; store != nil; store, err = store.NextPage() {
				for _, item := range store.Items {
					models.AddLPOffer(item.ID, npcCorp.ID, item.Item.ID, item.Quantity, item.LpCost, item.AkCost, item.IskCost)
					for _, requirement := range item.RequiredItems {
						models.AddLPOfferRequirements(item.ID, requirement.Item.ID, requirement.Quantity)
					}
				}
			}
		}
	}
	return nil
}
