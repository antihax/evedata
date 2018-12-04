package hammer

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
)

func init() {
	registerConsumer("characterStructureMarket", structureOrdersConsumer)
	registerConsumer("marketHistoryTrigger", marketHistoryTrigger)
	registerConsumer("marketHistory", marketHistoryConsumer)
}

func structureOrdersConsumer(s *Hammer, parameter interface{}) {
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))
	structureID := parameters[2].(int64)

	if s.inQueue.CheckWorkExpired("evedata_structuremarket_failure",
		fmt.Sprintf("%d%d", structureID, tokenCharacterID)) {
		return
	}

	// early out if we already have this recently
	if s.inQueue.CheckWorkExpired("evedata_structuremarket", structureID) {
		return
	}

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	var page int32 = 1
	orders := []esi.GetMarketsStructuresStructureId200Ok{}
	for {
		o, r, err := s.esi.ESI.MarketApi.GetMarketsStructuresStructureId(ctx, structureID,
			&esi.GetMarketsStructuresStructureIdOpts{
				Page: optional.NewInt32(page),
			})
		if err != nil {
			s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
			if r != nil && r.StatusCode == 403 {
				err := s.inQueue.SetWorkExpire("evedata_structuremarket_failure", fmt.Sprintf("%d%d", structureID, tokenCharacterID), 86400*3)
				if err != nil {
					log.Printf("failed setting failure: %s %d\n", err, structureID)
				}
			}
			return
		} else if len(o) == 0 { // end of the pages
			break
		}

		orders = append(orders, o...)

		page++
	}
	// early out if there are no orders
	if len(orders) == 0 {
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.StructureOrders{Orders: orders, StructureID: structureID}, "structureOrders")
	if err != nil {
		log.Println(err)
		return
	}
}

func marketHistoryTrigger(s *Hammer, parameter interface{}) {
	regions, _, err := s.esi.ESI.UniverseApi.GetUniverseRegions(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	var page int32 = 1

	for {
		items, r, err := s.esi.ESI.UniverseApi.GetUniverseTypes(context.Background(),
			&esi.GetUniverseTypesOpts{
				Page: optional.NewInt32(page),
			})
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
				s.QueueWork(work, redisqueue.Priority_Lowest)
			}
		}

		xpagesS := r.Header.Get("x-pages")
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
