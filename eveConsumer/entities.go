package eveConsumer

import (
	"errors"
	"fmt"
	"time"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("entities", characterConsumer, "EVEDATA_characterQueue")
	addConsumer("entities", allianceConsumer, "EVEDATA_allianceQueue")
	addConsumer("entities", corporationConsumer, "EVEDATA_corporationQueue")
	addConsumer("entities", charSearchConsumer, "EVEDATA_charSearchQueue")
	addTrigger("entities", entitiesTrigger)
}

// At the public rate limit, we can obtain 540,000 entities an hour.
// Recursion will be limited to once an day with expiration of entities at five days.

// Check if we need to update any entity information (character, corporation, alliance)
func entitiesTrigger(c *EVEConsumer) (bool, error) {

	r := c.ctx.Cache.Get()

	chars, err := models.MaintOrphanCharacters()
	if err != nil {
		return false, err
	}

	for _, char := range chars {
		r.Send("SADD", "EVEDATA_entityQueue", char)
	}
	r.Flush()
	r.Close()

	err = c.entitiesFromCREST()
	if err != nil {
		return false, err
	}

	err = c.allianceUpdate()
	if err != nil {
		return false, err
	}

	err = c.corporationUpdate()
	if err != nil {
		return false, err
	}

	err = c.characterUpdate()
	if err != nil {
		return false, err
	}

	return true, err
}

func charSearchConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
	ret, err := r.Do("SPOP", "EVEDATA_charSearchQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}
	v, err := redis.String(ret, err)
	if err != nil {
		return false, err
	}

	if !goesi.ValidCharacterName(v) {
		return false, errors.New(fmt.Sprintf("Invalid Character Name: %s", v))
	}

	// Figure out if we know this person already
	id, err := models.GetCharacterIDByName(v)
	if err != nil {
		return true, err
	}

	// We don't know this person... lets go looking.

	if id == 0 {
		search, _, err := c.ctx.ESI.ESI.SearchApi.GetSearch([]string{"character"}, v, map[string]interface{}{"strict": true})
		if err != nil {
			return true, err
		}
		if len(search.Character) > 0 {
			redis := c.ctx.Cache.Get()
			for _, nid := range search.Character {
				EntityCharacterAddToQueue(nid, &redis)
			}
			redis.Close()
		}
	} else { // add the character to the queue so we get latest data.
		redis := c.ctx.Cache.Get()
		EntityCharacterAddToQueue((int32)(id), &redis)
		redis.Close()
	}

	return true, err
}

func characterConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
	ret, err := r.Do("SPOP", "EVEDATA_characterQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}
	v, err := redis.Int(ret, err)
	if err != nil {
		return false, err
	}

	// Skip this entity if we have touched it recently
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", v)
	i, err := redis.Bool(r.Do("EXISTS", key))
	if err != nil || i == true {
		return false, err
	}

	err = c.entityGetAndSaveCategory((int32)(v), "character")
	if err != nil {
		return false, err
	}
	return true, err
}

func corporationConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
	ret, err := r.Do("SPOP", "EVEDATA_corporationQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}
	v, err := redis.Int(ret, err)
	if err != nil {
		return false, err
	}

	// Skip this entity if we have touched it recently
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", v)
	i, err := redis.Bool(r.Do("EXISTS", key))
	if err != nil || i == true {
		return false, err
	}

	err = c.entityGetAndSaveCategory((int32)(v), "corporation")
	if err != nil {
		return false, err
	}
	return true, err
}

func allianceConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
	ret, err := r.Do("SPOP", "EVEDATA_allianceQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}
	v, err := redis.Int(ret, err)
	if err != nil {
		return false, err
	}

	// Skip this entity if we have touched it recently
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", v)
	i, err := redis.Bool(r.Do("EXISTS", key))
	if err != nil || i == true {
		return false, err
	}

	err = c.entityGetAndSaveCategory((int32)(v), "alliance")
	if err != nil {
		return false, err
	}
	return true, err
}

// update any old entities
func (c *EVEConsumer) allianceUpdate() error {
	entities, err := c.ctx.Db.Query(
		`SELECT allianceid AS id, crestRef, cacheUntil FROM evedata.alliances A
			INNER JOIN evedata.crestID C ON A.allianceID = C.id
						WHERE cacheUntil < UTC_TIMESTAMP()
            ORDER BY cacheUntil ASC`)
	if err != nil {
		return err
	}

	r := c.ctx.Cache.Get()
	defer r.Close()

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
		if err = EntityAllianceAddToQueue(id, &r); err != nil {
			return err
		}

	}
	entities.Close()

	return nil
}
func (c *EVEConsumer) characterUpdate() error {
	entities, err := c.ctx.Db.Query(
		`SELECT characterID AS id, crestRef, cacheUntil FROM evedata.characters A
			INNER JOIN evedata.crestID C ON A.characterID = C.id
						WHERE cacheUntil < UTC_TIMESTAMP()  
            ORDER BY cacheUntil ASC`)
	if err != nil {
		return err
	}

	r := c.ctx.Cache.Get()
	defer r.Close()

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
		if err = EntityCharacterAddToQueue(id, &r); err != nil {
			return err
		}

	}
	entities.Close()

	return nil
}

func (c *EVEConsumer) corporationUpdate() error {
	entities, err := c.ctx.Db.Query(
		`SELECT corporationid AS id, crestRef, cacheUntil FROM evedata.corporations A
			INNER JOIN evedata.crestID C ON A.corporationID = C.id
						WHERE cacheUntil < UTC_TIMESTAMP() AND memberCount > 0
            ORDER BY cacheUntil ASC`)
	if err != nil {
		return err
	}

	r := c.ctx.Cache.Get()
	defer r.Close()

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
		if err = EntityCorporationAddToQueue(id, &r); err != nil {
			return err
		}

	}
	entities.Close()

	return nil
}

// Collect entity information for new alliances
func (c *EVEConsumer) entitiesFromCREST() error {

	nextCheck, _, err := models.GetServiceState("alliances")
	if err != nil {
		return err
	} else if nextCheck.After(time.Now().UTC()) {
		return nil
	}

	ids, res, err := c.ctx.ESI.ESI.AllianceApi.GetAlliances(nil)
	if err != nil {
		return err
	}

	// Update state so we dont have two polling at once.
	err = models.SetServiceState("alliances", goesi.CacheExpires(res).UTC(), 1)
	if err != nil {
		return err
	}

	redis := c.ctx.Cache.Get()
	defer redis.Close()
	// Throw them into the queue
	for _, allianceID := range ids {
		if err = EntityAllianceAddToQueue(allianceID, &redis); err != nil {
			return err
		}
	}

	return nil
}

func CharSearchAddToQueue(charList []interface{}, redisPtr *redis.Conn) {
	r := *redisPtr

	for _, name := range charList {
		if goesi.ValidCharacterName(name.(string)) {
			// Add the search to the queue
			r.Send("SADD", "EVEDATA_charSearchQueue", name.(string))
		}
	}
	r.Flush()
}

func EntityAllianceAddToQueue(id int32, r *redis.Conn) error {
	red := *r
	// Skip this entity if we have touched it recently
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", id)
	i, err := redis.Bool(red.Do("EXISTS", key))
	if err != nil || i == true {
		return err
	}

	// Add the entity to the queue
	_, err = red.Do("SADD", "EVEDATA_allianceQueue", id)
	return err
}

func EntityCorporationAddToQueue(id int32, r *redis.Conn) error {
	red := *r
	// Skip this entity if we have touched it recently
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", id)
	i, err := redis.Bool(red.Do("EXISTS", key))
	if err != nil || i == true {
		return err
	}

	// Add the entity to the queue
	_, err = red.Do("SADD", "EVEDATA_corporationQueue", id)
	return err
}

func EntityCharacterAddToQueue(id int32, r *redis.Conn) error {
	red := *r
	// Skip this entity if we have touched it recently
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", id)
	i, err := redis.Bool(red.Do("EXISTS", key))
	if err != nil || i == true {
		return err
	}

	// Add the entity to the queue
	_, err = red.Do("SADD", "EVEDATA_characterQueue", id)
	return err
}

// Say we touched the entity and expire after one day
func (c *EVEConsumer) entitySetKnown(id int32) error {
	key := "EVEDATA_entity:" + fmt.Sprintf("%d\n", id)
	r := c.ctx.Cache.Get()
	defer r.Close()
	r.Do("SETEX", key, 3600, true)
	return nil
}

func (c *EVEConsumer) entityGetAndSaveCategory(id int32, category string) error {
	var err error
	h := "https://crest-tq.eveonline.com/" + fmt.Sprintf("%ss/%d/", category, id)
	if category == "alliance" {
		err = c.updateAlliance(id)
	} else if category == "corporation" {
		err = c.updateCorporation(id)
	} else if category == "character" {
		err = c.updateCharacter(id)
	}
	if err != nil {
		return err
	}

	c.entitySetKnown(id)

	err = models.AddCRESTRef(((int64)(id)), h)
	if err != nil {
		return err
	}

	return err
}

func (c *EVEConsumer) updateAlliance(id int32) error {
	a, _, err := c.ctx.ESI.ESI.AllianceApi.GetAlliancesAllianceId(id, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("%s with alliance id %d", err, id))
	}

	corps, _, err := c.ctx.ESI.ESI.AllianceApi.GetAlliancesAllianceIdCorporations(id, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("%s with alliance id %d", err, id))
	}

	err = models.UpdateAlliance(id, a.AllianceName, len(corps), a.Ticker, a.ExecutorCorp,
		a.DateFounded, time.Now().UTC().Add(time.Hour*24))
	if err != nil {
		return errors.New(fmt.Sprintf("%s with alliance id %d", err, id))
	}

	redis := c.ctx.Cache.Get()
	defer redis.Close()
	for _, corp := range corps {
		err = EntityCorporationAddToQueue(corp, &redis)
		if err != nil {
			return errors.New(fmt.Sprintf("%s with alliance id %d", err, id))
		}
	}

	return nil
}

func (c *EVEConsumer) updateCorporation(id int32) error {
	a, _, err := c.ctx.ESI.ESI.CorporationApi.GetCorporationsCorporationId(id, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("%s with corporation id %d", err, id))
	}

	factionID := goesi.FactionNameToID(a.Faction)
	err = models.UpdateCorporation(id, a.CorporationName, a.Ticker, a.CeoId,
		a.AllianceId, factionID, a.MemberCount, time.Now().UTC().Add(time.Hour*24))
	if err != nil {
		return errors.New(fmt.Sprintf("%s with corporation id %d", err, id))
	}
	if a.CeoId > 1 {
		redis := c.ctx.Cache.Get()
		defer redis.Close()
		err = EntityCharacterAddToQueue((int32)(a.CeoId), &redis)
		if err != nil {
			return errors.New(fmt.Sprintf("%s with corporation id %d", err, id))
		}
	}
	return nil
}

func (c *EVEConsumer) updateCharacter(id int32) error {
	a, _, err := c.ctx.ESI.ESI.CharacterApi.GetCharactersCharacterId(id, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("%s with character id %d", err, id))
	}
	err = models.UpdateCharacter(id, a.Name, a.BloodlineId, a.AncestryId, a.CorporationId, a.AllianceId, a.RaceId, a.Gender, a.SecurityStatus, time.Now().UTC().Add(time.Hour*24))
	if err != nil {
		return errors.New(fmt.Sprintf("%s with character id %d", err, id))
	}

	redis := c.ctx.Cache.Get()
	defer redis.Close()
	err = EntityCharacterAddToQueue(id, &redis)

	h, _, err := c.ctx.ESI.ESI.CharacterApi.GetCharactersCharacterIdCorporationhistory(id, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("%s with character history id %d", err, id))
	}

	for _, corp := range h {
		err = models.UpdateCorporationHistory(id, corp.CorporationId, corp.RecordId, corp.StartDate)
		if err != nil {
			return errors.New(fmt.Sprintf("%s with character history id %d", err, id))
		}
	}
	return nil
}
