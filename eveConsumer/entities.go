package eveConsumer

import (
	"fmt"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("entities", charSearchConsumer, "EVEDATA_charSearchQueue")
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
		return false, fmt.Errorf("Invalid Character Name: %s", v)
	}

	// Figure out if we know this person already
	id, err := models.GetCharacterIDByName(v)
	if err != nil {
		return true, err
	}

	// We don't know this person... lets go looking.

	if id == 0 {
		search, _, err := c.ctx.ESI.ESI.SearchApi.GetSearch(nil, []string{"character"}, v, map[string]interface{}{"strict": true})
		if err != nil {
			return true, err
		}
		if len(search.Character) > 0 {
			redis := c.ctx.Cache.Get()
			for _, nid := range search.Character {
				EntityCharacterAddToQueue(nid)
			}
			redis.Close()
		}
	} else { // add the character to the queue so we get latest data.
		redis := c.ctx.Cache.Get()
		EntityCharacterAddToQueue((int32)(id))
		redis.Close()
	}

	return true, err
}

// CharSearchAddToQueue add a character to the search queue
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

func EntityAllianceAddToQueue(allianceID int32) error {
	if allianceID > 99000000 { // Skip NPC Alliances
		if !hammerQueue.CheckWorkExpired("evedata_entity", int64(allianceID)) {
			return hammerQueue.QueueWork([]redisqueue.Work{
				{Operation: "alliance", Parameter: allianceID},
			})
		}
	}
	return nil
}

func EntityCorporationAddToQueue(corporationID int32) error {
	if corporationID > 98000000 { // Skip NPC Corporations
		if !hammerQueue.CheckWorkExpired("evedata_entity", int64(corporationID)) {
			return hammerQueue.QueueWork([]redisqueue.Work{
				{Operation: "corporation", Parameter: corporationID},
			})
		}
	}
	return nil
}

func EntityCharacterAddToQueue(characterID int32) error {
	if characterID > 90000000 { // Skip NPC Characters
		if !hammerQueue.CheckWorkExpired("evedata_entity", int64(characterID)) {
			return hammerQueue.QueueWork([]redisqueue.Work{
				{Operation: "character", Parameter: characterID},
			})
		}
	}
	return nil
}
