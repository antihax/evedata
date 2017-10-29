package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"

	"encoding/gob"

	"github.com/antihax/evedata/internal/gobcoder"
)

func init() {
	registerConsumer("structure", structureConsumer)
	registerConsumer("structureOrders", structureOrdersConsumer)
	gob.Register(datapackages.Structure{})
	gob.Register(datapackages.StructureOrders{})
}

func structureConsumer(s *Hammer, parameter interface{}) {
	structureID := parameter.(int64)

	if s.inQueue.CheckWorkExpired("evedata_structure_failure", structureID) {
		log.Printf("ignoring structure %d\n", structureID)
		return
	}

	ctx := context.WithValue(context.TODO(), goesi.ContextOAuth2, s.token)
	struc, _, err := s.esi.ESI.UniverseApi.GetUniverseStructuresStructureId(ctx, structureID, nil)
	if err != nil {
		s.inQueue.SetWorkExpire("evedata_structure_failure", structureID, 86400)
		return
	}

	b, err := gobcoder.GobEncoder(datapackages.Structure{Structure: struc, StructureID: structureID})
	if err != nil {
		log.Println(err)
		return
	}
	err = s.nsq.Publish("structure", b)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func structureOrdersConsumer(s *Hammer, parameter interface{}) {
	structureID := parameter.(int64)
	var page int32 = 1
	orders := []esi.GetMarketsStructuresStructureId200Ok{}

	if s.inQueue.CheckWorkExpired("evedata_structure_failure", structureID) {
		log.Printf("ignoring structure %d\n", structureID)
		return
	}

	ctx := context.WithValue(context.TODO(), goesi.ContextOAuth2, s.token)

	for {
		o, _, err := s.esi.ESI.MarketApi.GetMarketsStructuresStructureId(ctx, structureID, map[string]interface{}{"page": page})
		if err != nil {
			s.inQueue.SetWorkExpire("evedata_structure_failure", structureID, 86400)
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

	b, err := gobcoder.GobEncoder(&datapackages.StructureOrders{Orders: orders, StructureID: structureID})
	if err != nil {
		log.Println(err)
		return
	}

	err = s.nsq.Publish("structureOrders", b)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
