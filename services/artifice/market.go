package artifice

import (
	"context"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	//registerTrigger("structures", structuresTrigger, time.NewTicker(time.Second*300))
	registerTrigger("marketHistory", historyTrigger, time.NewTicker(time.Second*3000))
}

func historyTrigger(s *Artifice) error {
	hour := time.Now().UTC().Hour()
	if hour == 1 {
		work := []redisqueue.Work{}
		work = append(work, redisqueue.Work{Operation: "marketHistoryTrigger", Parameter: true})
		return s.QueueWork(work, redisqueue.Priority_High)
	}
	return nil
}

func structuresTrigger(s *Artifice) error {
	// Get a list of structures to rake over for market data also
	structures, _, err := s.esi.ESI.UniverseApi.GetUniverseStructures(context.Background(), nil)
	if err != nil {
		return err
	}

	structure, err := s.inQueue.CheckWorkCompletedInBulk("evedata_structure_failure", structures)
	if err != nil {
		return err
	}

	for i := range structure {
		if !structure[i] {
			work := []redisqueue.Work{}
			work = append(work, redisqueue.Work{Operation: "structure", Parameter: structures[i]})
			err = s.QueueWork(work, redisqueue.Priority_Lowest)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
