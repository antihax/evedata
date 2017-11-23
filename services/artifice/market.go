package artifice

import (
	"context"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("marketOrders", marketTrigger, time.NewTicker(time.Second*300))
	registerTrigger("structures", structuresTrigger, time.NewTicker(time.Second*300))
	registerDailyTrigger("marketHistory", historyTrigger, 1)
}

func marketTrigger(s *Artifice) error {
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.Background(), nil)
	if err != nil {
		return err
	}

	work := []redisqueue.Work{}
	for _, region := range regions {
		if region < 11000000 {
			work = append(work, redisqueue.Work{Operation: "marketOrders", Parameter: region})
		}
	}
	s.QueueWork(work)
	return nil
}

func historyTrigger(s *Artifice) error {
	work := []redisqueue.Work{}
	work = append(work, redisqueue.Work{Operation: "marketHistoryTrigger", Parameter: true})
	s.QueueWork(work)
	return nil
}

func structuresTrigger(s *Artifice) error {
	// Get a list of structures to rake over for market data also
	structures, _, err := s.esi.ESI.UniverseApi.GetUniverseStructures(context.Background(), nil)
	if err != nil {
		return err
	}

	failed, err := s.inQueue.CheckWorkCompletedInBulk("evedata_structure_failure", structures)
	if err != nil {
		return err
	}

	work := []redisqueue.Work{}
	for i := range failed {
		if !failed[i] {
			work = append(work, redisqueue.Work{Operation: "structure", Parameter: structures[i]})
			work = append(work, redisqueue.Work{Operation: "structureOrders", Parameter: structures[i]})
		}
	}
	s.QueueWork(work)
	return nil
}
