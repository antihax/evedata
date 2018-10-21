package marketwatch

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *MarketWatch) contractWorker(regionID int32) {
	// For totalization
	wg := sync.WaitGroup{}

	// Loop forever
	for {
		start := time.Now()
		numContracts := 0

		// Return Channels
		rchan := make(chan []esi.GetContractsPublicRegionId200Ok, 100000)
		echan := make(chan error, 100000)

		contracts, res, err := s.esi.ESI.ContractsApi.GetContractsPublicRegionId(
			context.Background(), regionID, nil,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		rchan <- contracts

		// Figure out if there are more pages
		pages, err := getPages(res)
		if err != nil {
			log.Println(err)
			continue
		}
		duration := timeUntilCacheExpires(res)
		if duration.Minutes() < 3 {
			fmt.Printf("%d contract too close to window: waiting %s\n", regionID, duration.String())
			time.Sleep(duration)
			continue
		}

		// Get the other pages concurrently
		for pages > 1 {
			wg.Add(1) // count whats running
			go func(page int32) {
				defer wg.Done() // release when done

				contracts, _, err := s.esi.ESI.ContractsApi.GetContractsPublicRegionId(
					context.Background(),
					regionID,
					&esi.GetContractsPublicRegionIdOpts{Page: optional.NewInt32(page)},
				)
				if err != nil {
					echan <- err
					return
				}

				// Add the contracts to the channel
				rchan <- contracts
			}(pages)
			pages--
		}

		wg.Wait() // Wait for everything to finish

		// Close the channels
		close(rchan)
		close(echan)

		for err := range echan {
			log.Println(err)
		}

		changes := []ContractChange{}
		newContracts := []FullContract{}
		// Add all the contracts together
		for o := range rchan {
		Restart:
			for i := range o {

				// Ignore expired contracts
				if o[i].DateExpired.Before(time.Now()) {
					continue
				}

				contract := Contract{Touched: start, Contract: FullContract{Contract: o[i]}}

				if o[i].Type_ == "item_exchange" || o[i].Type_ == "auction" {
					err := s.getContractItems(&contract)
					if err != nil {
						goto Restart
					}
				}

				if o[i].Type_ == "auction" {
					err := s.getContractBids(&contract)
					if err != nil {
						goto Restart
					}
				}

				change, isNew := s.storeContract(int64(regionID), contract)
				numContracts++
				if change.Changed && !isNew {
					changes = append(changes, change)
				}
				if isNew {
					newContracts = append(newContracts, contract.Contract)
				}
			}
		}
		deletions := s.expireContracts(int64(regionID), start)

		// Log metrics
		metricContractTimePull.With(
			prometheus.Labels{
				"locationID": strconv.FormatInt(int64(regionID), 10),
			},
		).Observe(float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond))

		if len(newContracts) > 0 {
			s.contractChan <- newContracts
		}

		// Only bids really change.
		if len(changes) > 0 {
			s.contractChangeChan <- changes
		}

		if len(deletions) > 0 {
			s.contractDeleteChan <- changes
		}

		// Sleep until the cache timer expires, plus a little.
		time.Sleep(duration)
	}
}

// getContractItems for a single contract. Must be prefilled with the contract.
func (s *MarketWatch) getContractItems(contract *Contract) error {
	wg := sync.WaitGroup{}

	// Return Channels
	rchan := make(chan []esi.GetContractsPublicItemsContractId200Ok, 100000)
	echan := make(chan error, 100000)

	items, res, err := s.esi.ESI.ContractsApi.GetContractsPublicItemsContractId(
		context.Background(), contract.Contract.Contract.ContractId, nil,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	// No items on the order
	if res.StatusCode == 204 || res.StatusCode == 403 {
		return nil
	}

	rchan <- items

	// Figure out if there are more pages
	pages, err := getPages(res)
	if err != nil {
		log.Printf("%d %v\n", contract.Contract.Contract.ContractId, err)

		return err
	}

	// Get the other pages concurrently
	for pages > 1 {
		wg.Add(1) // count whats running
		go func(page int32) {
			defer wg.Done() // release when done

			items, _, err := s.esi.ESI.ContractsApi.GetContractsPublicItemsContractId(
				context.Background(),
				contract.Contract.Contract.ContractId,
				&esi.GetContractsPublicItemsContractIdOpts{Page: optional.NewInt32(page)},
			)
			if err != nil {
				echan <- err
				return
			}

			// Add the contracts to the channel
			rchan <- items
		}(pages)
		pages--
	}

	wg.Wait()

	// Close the channels
	close(rchan)
	close(echan)

	for err := range echan {
		// Fail all if one fails
		log.Println(err)
		return err
	}

	// Add all the contracts together
	for o := range rchan {
		contract.Contract.Items = append(contract.Contract.Items, o...)
	}

	return nil
}

// getContractBids for a single contract. Must be prefilled with the contract.
func (s *MarketWatch) getContractBids(contract *Contract) error {
	wg := sync.WaitGroup{}

	// Return Channels
	rchan := make(chan []esi.GetContractsPublicBidsContractId200Ok, 100000)
	echan := make(chan error, 100000)

	bids, res, err := s.esi.ESI.ContractsApi.GetContractsPublicBidsContractId(
		context.Background(), contract.Contract.Contract.ContractId, nil,
	)
	rchan <- bids

	// Figure out if there are more pages
	pages, err := getPages(res)
	if err != nil {
		log.Println(err)
		return err
	}

	// Get the other pages concurrently
	for pages > 1 {
		wg.Add(1) // count whats running
		go func(page int32) {
			defer wg.Done() // release when done

			bids, _, err := s.esi.ESI.ContractsApi.GetContractsPublicBidsContractId(
				context.Background(),
				contract.Contract.Contract.ContractId,
				&esi.GetContractsPublicBidsContractIdOpts{Page: optional.NewInt32(page)},
			)
			if err != nil {
				echan <- err
				return
			}

			// Add the contracts to the channel
			rchan <- bids
		}(pages)
		pages--
	}

	wg.Wait()

	// Close the channels
	close(rchan)
	close(echan)

	for err := range echan {
		// Fail all if one fails
		log.Println(err)
		return err
	}

	// Add all the bids together
	for o := range rchan {
		contract.Contract.Bids = append(contract.Contract.Bids, o...)
	}

	return nil
}

// Metrics
var (
	metricContractTimePull = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "contract",
		Name:      "pull",
		Help:      "Market Pull Statistics",
		Buckets:   prometheus.ExponentialBuckets(10, 1.6, 20),
	},
		[]string{"locationID"},
	)
)

func init() {
	prometheus.MustRegister(
		metricContractTimePull,
	)
}
