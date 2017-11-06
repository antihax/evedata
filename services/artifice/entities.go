package artifice

import (
	"context"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("alliance", allianceTrigger, time.NewTicker(time.Second*3600))
	registerTrigger("characterUpdate", characterUpdate, time.NewTicker(time.Second*120))
	registerTrigger("corporationUpdate", corporationUpdate, time.NewTicker(time.Second*120))
}

func allianceTrigger(s *Artifice) error {
	alliances, _, err := s.esi.ESI.AllianceApi.GetAlliances(context.TODO(), nil)
	if err != nil {
		return err
	}

	work := []redisqueue.Work{}
	for _, alliance := range alliances {
		work = append(work, redisqueue.Work{Operation: "alliance", Parameter: alliance})
	}
	s.QueueWork(work)
	return nil
}

func characterUpdate(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT characterID AS id FROM evedata.characters A
			WHERE cacheUntil < UTC_TIMESTAMP() AND dead = 0
			AND characterID > 90000000
            ORDER BY cacheUntil ASC`)
	if err != nil {
		return err
	}

	work := []redisqueue.Work{}

	// Loop the entities
	for entities.Next() {
		var id int32

		err = entities.Scan(&id)
		if err != nil {
			return err
		}

		work = append(work, redisqueue.Work{Operation: "character", Parameter: id})

	}
	s.QueueWork(work)
	entities.Close()

	return nil
}

func corporationUpdate(s *Artifice) error {
	entities, err := s.db.Query(
		`SELECT corporationID AS id FROM evedata.corporations A
		 WHERE cacheUntil < UTC_TIMESTAMP() AND memberCount > 0 AND corporationId> 90000000
		 ORDER BY cacheUntil ASC`)
	if err != nil {
		return err
	}

	work := []redisqueue.Work{}

	// Loop the entities
	for entities.Next() {
		var id int32

		err = entities.Scan(&id)
		if err != nil {
			return err
		}

		work = append(work, redisqueue.Work{Operation: "corporation", Parameter: id})

	}
	s.QueueWork(work)
	entities.Close()

	return nil
}
