package hammer

import (
	"log"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/goesi/esi"

	"encoding/gob"

	"github.com/antihax/evedata/internal/gobcoder"
)

func init() {
	registerConsumer("marketOrders", marketOrdersConsumer)
	gob.Register(datapackages.MarketOrders{})
}

func marketOrdersConsumer(s *Hammer, parameter interface{}) {
	regionID := parameter.(int32)
	var page int32 = 1
	orders := []esi.GetMarketsRegionIdOrders200Ok{}

	for {
		o, _, err := s.esi.ESI.MarketApi.GetMarketsRegionIdOrders(nil, "all", regionID, map[string]interface{}{"page": page})
		if err != nil {
			log.Println(err)
			return
		} else if len(o) == 0 { // end of the pages
			break
		}
		orders = append(orders, o...)

		page++
	}

	b, err := gobcoder.GobEncoder(&datapackages.MarketOrders{Orders: orders, RegionID: regionID})
	if err != nil {
		log.Println(err)
		return
	}
	err = s.nsq.Publish("marketOrders", b)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
