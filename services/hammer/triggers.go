package hammer

import (
	"log"
	"time"

	"github.com/antihax/evedata/internal/redisqueue"
)

func (s *Hammer) runTriggers() {
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	marketOrders := time.NewTicker(20 * time.Minute)

	go func() {
		for {
			select {
			case <-marketOrders.C:
				work := []redisqueue.Work{}
				for _, region := range regions {
					work = append(work, redisqueue.Work{Operation: "marketOrders", Parameter: region})
				}

				s.QueueWork(work)
			case <-s.stop:
				marketOrders.Stop()
				return
			}
		}
	}()

}
