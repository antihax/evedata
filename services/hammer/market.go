package hammer

import (
	"context"
	"log"
	"strconv"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/redisqueue"

	"encoding/gob"
)

func init() {
	registerConsumer("marketOrders", marketOrdersConsumer)
	registerConsumer("marketHistoryTrigger", marketHistoryTrigger)
	registerConsumer("marketHistory", marketHistoryConsumer)
	gob.Register(datapackages.MarketOrders{})
	gob.Register(datapackages.MarketHistory{})
}

func marketHistoryTrigger(s *Hammer, parameter interface{}) {
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	var page int32 = 1

	for {
		items, r, err := s.esi.ESI.UniverseApi.GetUniverseTypes(context.TODO(), map[string]interface{}{"page": page})
		if err != nil {
			log.Println(err)
			continue
		}

		for _, itemID := range items {
			item, _, err := s.esi.ESI.UniverseApi.GetUniverseTypesTypeId(context.TODO(), itemID, nil)
			if err != nil {
				continue
			}
			if item.Published && item.MarketGroupId > 0 {
				work := []redisqueue.Work{}
				for _, regionID := range regions {
					if regionID < 11000000 {
						work = append(work, redisqueue.Work{Operation: "marketHistory", Parameter: []int32{regionID, itemID}})
					}
				}
				s.QueueWork(work)
			}
		}

		xpagesS := r.Header.Get("X-Pages")
		xpages, _ := strconv.Atoi(xpagesS)
		if int32(xpages) == page || len(items) == 0 {
			return
		}
		page++
	}
}

func marketHistoryConsumer(s *Hammer, parameter interface{}) {
	regionID := parameter.([]int32)[0]
	typeID := parameter.([]int32)[1]
	h, _, err := s.esi.ESI.MarketApi.GetMarketsRegionIdHistory(nil, regionID, typeID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.MarketHistory{
		History:  h,
		RegionID: regionID,
		TypeID:   typeID}, "marketHistory")
	if err != nil {
		log.Println(err)
		return
	}
}

func marketOrdersConsumer(s *Hammer, parameter interface{}) {
	regionID := parameter.(int32)
	var page int32 = 1

	for {
		orders, r, err := s.esi.ESI.MarketApi.GetMarketsRegionIdOrders(context.TODO(), "all", regionID, map[string]interface{}{"page": page})
		if err != nil {
			log.Println(err)
			return
		}

		// Send out the result
		err = s.QueueResult(&datapackages.MarketOrders{
			Orders:   orders,
			RegionID: regionID},
			"marketOrders")
		if err != nil {
			log.Println(err)
			return
		}

		xpagesS := r.Header.Get("X-Pages")
		xpages, _ := strconv.Atoi(xpagesS)
		if int32(xpages) == page || len(orders) == 0 {
			return
		}
		page++
	}
}
