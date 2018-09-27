package marketwatch

import (
	"context"
	"log"
	"time"
)

func (s *MarketWatch) startUpMarketWorkers() {
	// Get all the regions and fire up workers for each
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, region := range regions {
		// Prebuild the maps
		s.createMarketStore(int64(region))
		s.createContractStore(int64(region))
		// Ignore non-market regions
		if region < 11000000 || region == 11000031 {
			time.Sleep(time.Millisecond * 500)
			go s.marketWorker(region)
			go s.contractWorker(region)
		}
	}

	if s.doAuth {
		go s.runStructures()
	}
}
