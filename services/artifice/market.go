package artifice

import (
	"context"
	"log"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("marketOrders", marketTrigger, time.NewTicker(time.Minute*30))
	registerTrigger("structures", structuresTrigger, time.NewTicker(time.Minute*60))
	registerDailyTrigger("marketHistory", historyTrigger, 1)
}

func marketTrigger(s *Artifice) error {
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.TODO(), nil)
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
	structures, _, err := s.esi.ESI.UniverseApi.GetUniverseStructures(context.TODO(), nil)
	if err != nil {
		return err
	}

	work := []redisqueue.Work{}
	for _, structure := range structures {
		if !s.inQueue.CheckWorkExpired("evedata_structure_failure", structure) {
			work = append(work, redisqueue.Work{Operation: "structure", Parameter: structure})
			work = append(work, redisqueue.Work{Operation: "structureOrders", Parameter: structure})
		} else {
			log.Printf("wont queue structure %d\n", structure)
		}
	}
	s.QueueWork(work)
	return nil
}
