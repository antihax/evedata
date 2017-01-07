package eveConsumer

import (
	"fmt"
	"log"
	"time"

	"github.com/antihax/evedata/esi"
	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"
)

// At the public rate limit, we can obtain 540,000 entities an hour.
// Recursion will be limited to once an day with expiration of entities at five days.

// Check if we need to update any entity information (character, corporation, alliance)
func (c *EVEConsumer) checkEntities() {
	err := c.collectEntitiesFromCREST()
	if err != nil {
		log.Printf("EVEConsumer: collecting entities: %v", err)
	}
	err = c.updateEntities()
	if err != nil {
		log.Printf("EVEConsumer: updating entities: %v", err)
	}

}

// update any old entities
func (c *EVEConsumer) updateEntities() error {

	entities, err := c.ctx.Db.Query(
		`SELECT allianceid AS id, crestRef, cacheUntil FROM evedata.alliances A
			INNER JOIN crestID C ON A.allianceID = C.id
						WHERE cacheUntil < UTC_TIMESTAMP()  
			UNION
			SELECT corporationid AS id, crestRef, cacheUntil FROM evedata.corporations A
			INNER JOIN crestID C ON A.corporationID = C.id
						WHERE cacheUntil < UTC_TIMESTAMP()
			UNION
			(SELECT characterID AS id, crestRef, cacheUntil FROM evedata.characters A
			INNER JOIN crestID C ON A.characterID = C.id
						WHERE cacheUntil < UTC_TIMESTAMP())
            
            ORDER BY cacheUntil ASC`)
	if err != nil {
		return err
	}

	// Loop the entities
	for entities.Next() {
		var (
			id      int32
			href    string
			nothing string
		)

		err = entities.Scan(&id, &href, &nothing)
		if err != nil {
			return err
		}

		// Recursively update expired information
		if err = c.entityAddToQueue(id); err != nil {
			return err
		}

	}
	entities.Close()

	return nil
}

// Collect entity information for new alliances
func (c *EVEConsumer) collectEntitiesFromCREST() error {

	nextCheck, _, err := models.GetServiceState("alliances")
	if err != nil {
		return err
	} else if nextCheck.Before(time.Now()) {
		return nil
	}

	// Get first page of alliances
	w, err := c.ctx.EVE.AlliancesV2(1)
	if err != nil {
		return err
	}

	// Update state so we dont have two polling at once.
	err = models.SetServiceState("alliances", w.CacheUntil, 1)
	if err != nil {
		return err
	}

	// Loop through all of the alliance pages
	for ; w != nil; w, err = w.NextPage() {
		if err != nil {
			return err
		}

		// Recursively update expired information
		for _, r := range w.Items {
			if err = c.entityAddToQueue((int32)(r.ID)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EVEConsumer) entityCheckQueue(r redis.Conn) error {
	ret, err := r.Do("SPOP", "EVEDATA_entityQueue")
	if err != nil {
		return err
	} else if ret == nil {
		return nil
	}
	v, err := redis.Int(ret, err)
	if err != nil {
		return err
	}

	// Skip this entity if we have touched it recently
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", v)
	i, err := redis.Bool(r.Do("EXISTS", key))
	if err != nil || i == true {
		return err
	}

	err = c.entityGetAndSave((int32)(v))
	return err
}

func (c *EVEConsumer) entityAddToQueue(id int32) error {
	r := c.ctx.Cache.Get()
	defer r.Close()

	// Skip this entity if we have touched it recently
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", id)
	i, err := redis.Bool(r.Do("EXISTS", key))
	if err != nil || i == true {
		return err
	}

	// Add the entity to the queue
	_, err = r.Do("SADD", "EVEDATA_entityQueue", id)
	return err
}

// Say we touched the entity and expire after one day
func (c *EVEConsumer) entitySetKnown(id int32) error {
	go func() {
		key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", id)
		r := c.ctx.Cache.Get()
		defer r.Close()
		r.Do("SETEX", key, 86400, true)
	}()
	return nil
}

// [TODO] Rewrite this as ESI matures
// [TODO] bulk pull IDs
func (c *EVEConsumer) entityGetAndSave(id int32) error {

	entity, _, err := c.ctx.ESI.UniverseApi.PostUniverseNames(esi.PostUniverseNamesIds{Ids: []int32{id}}, nil)
	if err != nil {
		return err
	}

	for _, e := range entity {
		h := "https://crest-tq.eveonline.com/" + fmt.Sprintf("%ss/%d/", e.Category, id)
		if e.Category == "alliance" {
			err = c.updateAlliance((int64)(e.Id))
		} else if e.Category == "corporation" {
			err = c.updateCorporation((int64)(e.Id))
		} else if e.Category == "character" {
			err = c.updateCharacter((int64)(e.Id))
		}

		if err != nil {
			return err
		}
		err = models.AddCRESTRef(((int64)(e.Id)), h)
		if err != nil {
			return err
		}
	}
	return err
}

func (c *EVEConsumer) updateAlliance(id int64) error {
	href := "https://crest-tq.eveonline.com/" + fmt.Sprintf("alliances/%d/", id)
	a, err := c.ctx.EVE.Alliance(href)
	if err != nil {
		return err
	}

	err = models.UpdateAlliance(a.ID, a.Name, a.CorporationsCount, a.ShortName, a.ExecutorCorporation.ID,
		a.StartDate.UTC(), a.Deleted, a.Description, a.CreatorCorporation.ID, a.CreatorCharacter.ID, time.Now().UTC().Add(time.Hour*24))
	if err != nil {
		return err
	}
	err = c.entityAddToQueue((int32)(a.CreatorCharacter.ID))
	if err != nil {
		return err
	}

	for _, corp := range a.Corporations {
		err = c.entityAddToQueue((int32)(corp.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *EVEConsumer) updateCorporation(id int64) error {
	a, err := c.ctx.EVE.CorporationPublicSheetXML(id)
	if err != nil {
		return err
	}

	err = models.UpdateCorporation(a.CorporationID, a.CorporationName, a.Ticker, a.CEOID, a.StationID,
		a.Description, a.AllianceID, a.FactionID, a.URL, a.MemberCount, a.Shares, time.Now().UTC().Add(time.Hour*24))
	if err != nil {
		return err
	}
	if a.CEOID > 1 {
		err = c.entityAddToQueue((int32)(a.CEOID))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *EVEConsumer) updateCharacter(id int64) error {
	a, err := c.ctx.EVE.CharacterInfoXML(id)
	if err != nil {
		return err
	}
	err = models.UpdateCharacter(a.CharacterID, a.CharacterName, a.BloodlineID, a.AncestryID, a.CorporationID, a.AllianceID, a.Race, a.SecurityStatus, time.Now().UTC().Add(time.Hour*24))
	if err != nil {
		return err
	}

	return nil
}
