package eveConsumer

import (
	"log"
	"time"

	"github.com/antihax/evedata/models"
)

func init() {
	addTrigger("npcCorps", npcCorpTrigger)
}

func npcCorpTrigger(c *EVEConsumer) error {
	nextCheck, _, err := models.GetServiceState("npcCorps")
	if err != nil {
		return err
	} else if nextCheck.After(time.Now()) {
		return nil
	}

	log.Printf("EVEConsumer: collecting loyalty Point Store Items")
	w, err := c.ctx.EVE.NPCCorporationsV1(1)
	if err != nil {
		return err
	}

	// Update state so we dont have two polling at once.
	err = models.SetServiceState("npcCorps", w.CacheUntil, 1)
	if err != nil {
		return err
	}

	redis := c.ctx.Cache.Get()
	defer redis.Close()

	// Loop through all of the pages
	for ; w != nil; w, err = w.NextPage() {
		for _, npcCorp := range w.Items {
			if npcCorp.LoyaltyStore.Href == "" {
				continue
			}

			EntityAddToQueue((int32)(npcCorp.ID), redis)
			store, err := c.ctx.EVE.LoyaltyPointStoreV1(npcCorp.LoyaltyStore.Href)
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
	return err
}
