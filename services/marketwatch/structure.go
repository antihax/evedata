package marketwatch

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/prometheus/client_golang/prometheus"
)

type Structure struct {
	restart time.Time
	running bool
}

func (s *MarketWatch) getAuthContext() context.Context {
	return context.WithValue(context.Background(), goesi.ContextOAuth2, *s.token)
}

func (s *MarketWatch) runStructures() {
	for {
		// Get all the structures and fire up workers for each
		structures, res, err := s.esi.ESI.UniverseApi.GetUniverseStructures(s.getAuthContext(), nil)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, structure := range structures {
			state := s.getStructureState(structure)
			// Prebuild the maps
			if state == nil {
				s.createMarketStore(structure)
				state = s.createStructureState(structure)
			}
			if state.running == false && time.Now().After(state.restart) {
				time.Sleep(time.Second * 1)
				state.running = true
				go s.structureWorker(structure)
			}
		}
		time.Sleep(timeUntilCacheExpires(res))
	}
}

func (s *MarketWatch) failStructure(structureID int64) {
	state := s.getStructureState(structureID)
	state.restart = time.Now().Add(time.Hour * 12)
	state.running = false
}

func (s *MarketWatch) structureWorker(structureID int64) {
	// For totalization
	wg := sync.WaitGroup{}

	// Loop forever
	for {
		start := time.Now()
		numOrders := 0

		// Return Channels
		rchan := make(chan []esi.GetMarketsStructuresStructureId200Ok, 100000)
		echan := make(chan error, 100000)

		orders, res, err := s.esi.ESI.MarketApi.GetMarketsStructuresStructureId(
			s.getAuthContext(), structureID, nil,
		)
		if err != nil {
			// If we do not have access, get out of the loop.
			if err.Error() == "403 Forbidden" {
				s.failStructure(structureID)
				return
			}
			log.Println(err)
			continue
		}

		rchan <- orders
		duration := timeUntilCacheExpires(res)
		if duration.Minutes() < 3 {
			fmt.Printf("%d too close to window: waiting %s\n", structureID, duration.String())
			time.Sleep(duration)
			continue
		}
		// Figure out if there are more pages
		pages, err := getPages(res)
		if err != nil {
			log.Println(err)
			continue
		}

		// Get the other pages concurrently
		for pages > 1 {
			wg.Add(1) // count whats running
			go func(page int32) {
				defer wg.Done() // release when done

				orders, r, err := s.esi.ESI.MarketApi.GetMarketsStructuresStructureId(
					context.Background(),
					structureID,
					&esi.GetMarketsStructuresStructureIdOpts{Page: optional.NewInt32(page)},
				)

				if err != nil {
					echan <- err
					return
				}

				// Are we too close to the end of the window?
				duration := timeUntilCacheExpires(r)
				if duration.Seconds() < 20 {
					echan <- errors.New("too close to end of window")
					return
				}

				// Add the orders to the channel
				rchan <- orders
			}(pages)
			pages--
		}

		wg.Wait() // Wait for everything to finish

		// Close the channels
		close(rchan)
		close(echan)

		for err := range echan {
			// Start over if any requests failed
			log.Println(err)
			continue
		}

		changes := []OrderChange{}
		newOrders := []esi.GetMarketsRegionIdOrders200Ok{}
		// Add all the orders together
		for o := range rchan {
			for i := range o {
				change, isNew := s.storeData(structureID,
					Order{
						Touched: start,
						Order:   sToR(o[i]),
					})
				numOrders++
				if change.Changed && !isNew {
					changes = append(changes, change)
				}
				if isNew {
					newOrders = append(newOrders, sToR(o[i]))
				}
			}
		}
		deletions := s.expireOrders(structureID, start)

		// Log metrics
		metricMarketTimePull.With(
			prometheus.Labels{
				"locationID": strconv.FormatInt(structureID, 10),
			},
		).Observe(float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond))

		if len(newOrders) > 0 {
			s.orderChan <- newOrders
		}

		/*if len(changes) > 0 {
			s.orderChangeChan <- changes
		}*/

		if len(deletions) > 0 {
			s.orderDeleteChan <- deletions
		}

		// Sleep until the cache timer expires
		time.Sleep(duration)
	}
}

// helper to copy a structure to a region to simplify storage.
func sToR(o esi.GetMarketsStructuresStructureId200Ok) esi.GetMarketsRegionIdOrders200Ok {
	return esi.GetMarketsRegionIdOrders200Ok{
		Duration:     o.Duration,
		IsBuyOrder:   o.IsBuyOrder,
		Issued:       o.Issued,
		LocationId:   o.LocationId,
		MinVolume:    o.MinVolume,
		OrderId:      o.OrderId,
		Price:        o.Price,
		Range_:       o.Range_,
		TypeId:       o.TypeId,
		VolumeRemain: o.VolumeRemain,
		VolumeTotal:  o.VolumeTotal,
	}
}
