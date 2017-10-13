package artifice

import (
	"context"

	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerTrigger("marketOrders", marketTrigger, 30)
	registerTrigger("structures", structuresTrigger, 60)
}

func marketTrigger(s *Artifice) error {
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.TODO(), nil)
	if err != nil {
		return err
	}

	work := []redisqueue.Work{}
	for _, region := range regions {
		work = append(work, redisqueue.Work{Operation: "marketOrders", Parameter: region})
	}
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
		work = append(work, redisqueue.Work{Operation: "structure", Parameter: structure})
		work = append(work, redisqueue.Work{Operation: "structureOrders", Parameter: structure})
	}
	s.QueueWork(work)
	return nil
}
