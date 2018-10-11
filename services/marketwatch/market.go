package marketwatch

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *MarketWatch) marketWorker(regionID int32) {
	// For totalization
	wg := sync.WaitGroup{}

	// Loop forever
	for {
		start := time.Now()
		numOrders := 0

		// Return Channels
		rchan := make(chan []esi.GetMarketsRegionIdOrders200Ok, 100000)
		echan := make(chan error, 100000)

		orders, res, err := s.esi.ESI.MarketApi.GetMarketsRegionIdOrders(
			context.Background(), "all", regionID, nil,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		rchan <- orders

		// Figure out if there are more pages
		pages, err := getPages(res)
		if err != nil {
			log.Println(err)
			continue
		}
		duration := timeUntilCacheExpires(res)
		if duration.Minutes() < 3 {
			fmt.Printf("%d market too close to window: waiting %s\n", regionID, duration.String())
			time.Sleep(duration)
			continue
		}

		// Get the other pages concurrently
		for pages > 1 {
			wg.Add(1) // count whats running
			go func(page int32) {
				defer wg.Done() // release when done

				orders, r, err := s.esi.ESI.MarketApi.GetMarketsRegionIdOrders(
					context.Background(),
					"all",
					regionID,
					&esi.GetMarketsRegionIdOrdersOpts{Page: optional.NewInt32(page)},
				)
				if err != nil {
					echan <- err
					return
				}

				// Are we too close to the end of the window?
				duration := timeUntilCacheExpires(r)
				if duration.Seconds() < 20 {
					echan <- errors.New("market too close to end of window")
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
				change, isNew := s.storeData(int64(regionID), Order{Touched: start, Order: o[i]})
				numOrders++
				if change.Changed && !isNew {
					changes = append(changes, change)
				}
				if isNew {
					newOrders = append(newOrders, o[i])
				}
			}
		}
		deletions := s.expireOrders(int64(regionID), start)

		// Log metrics
		metricMarketTimePull.With(
			prometheus.Labels{
				"locationID": strconv.FormatInt(int64(regionID), 10),
			},
		).Observe(float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond))

		//	fmt.Printf("%s %d (%d)\n", time.Since(start).String(), regionID, numOrders)

		if len(newOrders) > 0 {
			s.orderChan <- newOrders
		}

		/*	if len(changes) > 0 {
			s.orderChangeChan <- changes
		}*/

		if len(deletions) > 0 {
			s.orderDeleteChan <- deletions
		}

		// Sleep until the cache timer expires, plus a little.
		time.Sleep(duration)
	}
}

// Metrics
var (
	metricMarketTimePull = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "market",
		Name:      "pull",
		Help:      "Market Pull Statistics",
		Buckets:   prometheus.ExponentialBuckets(10, 1.6, 20),
	},
		[]string{"locationID"},
	)
)

func init() {
	prometheus.MustRegister(
		metricMarketTimePull,
	)
}
