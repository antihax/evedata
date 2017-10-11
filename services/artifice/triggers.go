package artifice

import (
	"context"
	"log"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func (s *Artifice) runTriggers() {
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	marketOrders := time.NewTicker(20 * time.Minute)

	go func() {
		for {
			select {
			// Get market orders every 20 minutes
			case <-marketOrders.C:
				work := []redisqueue.Work{}
				for _, region := range regions {
					work = append(work, redisqueue.Work{Operation: "marketOrders", Parameter: region})
				}
				s.QueueWork(work)

				work = []redisqueue.Work{}
				// Get a list of structures to rake over for market data also
				structures, _, err := s.esi.ESI.UniverseApi.GetUniverseStructures(context.TODO(), nil)
				if err != nil {
					log.Panicln(err)
					continue
				}
				for _, structure := range structures {
					work = append(work, redisqueue.Work{Operation: "structureOrders", Parameter: structure})
					work = append(work, redisqueue.Work{Operation: "structure", Parameter: structure})
				}
				s.QueueWork(work)

			case <-s.stop:
				marketOrders.Stop()
				return
			}
		}
	}()

}
