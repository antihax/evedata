package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"

	"encoding/gob"
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
		return
	}

	ctx := context.WithValue(context.TODO(), goesi.ContextOAuth2, s.token)
	struc, _, err := s.esi.ESI.UniverseApi.GetUniverseStructuresStructureId(ctx, structureID, nil)
	if err != nil {
		s.inQueue.SetWorkExpire("evedata_structure_failure", structureID, 86400)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.Structure{Structure: struc, StructureID: structureID}, "structure")
	if err != nil {
		log.Println(err)
		return
	}
}

func structureOrdersConsumer(s *Hammer, parameter interface{}) {
	structureID := parameter.(int64)
	var page int32 = 1
	orders := []esi.GetMarketsStructuresStructureId200Ok{}

	if s.inQueue.CheckWorkExpired("evedata_structure_failure", structureID) {
		return
	}

	ctx := context.WithValue(context.TODO(), goesi.ContextOAuth2, s.token)

	for {
		o, _, err := s.esi.ESI.MarketApi.GetMarketsStructuresStructureId(ctx, structureID, map[string]interface{}{"page": page})
		if err != nil {
			s.inQueue.SetWorkExpire("evedata_structure_failure", structureID, 172800)
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
	err := s.QueueResult(&datapackages.StructureOrders{Orders: orders, StructureID: structureID}, "structureOrders")
	if err != nil {
		log.Println(err)
		return
	}
}
