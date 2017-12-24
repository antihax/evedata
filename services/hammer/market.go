package hammer

import (
	"context"
	"log"
	"strconv"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/redisqueue"
)

func init() {
	registerConsumer("marketOrders", marketOrdersConsumer)
	registerConsumer("marketHistoryTrigger", marketHistoryTrigger)
	registerConsumer("marketHistory", marketHistoryConsumer)
}

func marketHistoryTrigger(s *Hammer, parameter interface{}) {
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	var page int32 = 1

	for {
		items, r, err := s.esi.ESI.UniverseApi.GetUniverseTypes(context.Background(), map[string]interface{}{"page": page})
		if err != nil {
			log.Println(err)
			continue
		}

		for _, itemID := range items {
			item, _, err := s.esi.ESI.UniverseApi.GetUniverseTypesTypeId(context.Background(), itemID, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			if item.Published && item.MarketGroupId > 0 {
				work := []redisqueue.Work{}
				for _, regionID := range regions {
					work = append(work, redisqueue.Work{Operation: "marketHistory", Parameter: []int32{regionID, itemID}})
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
	regionID := int32(parameter.([]interface{})[0].(int))
	typeID := int32(parameter.([]interface{})[1].(int))
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
	regionID := int32(parameter.(int))
	var page int32 = 1

	for {
		orders, r, err := s.esi.ESI.MarketApi.GetMarketsRegionIdOrders(context.Background(), "all", regionID, map[string]interface{}{"page": page})
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
