package hammer

import (
	"log"

	"encoding/gob"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/v1"
)

func init() {
	registerConsumer("marketOrders", marketOrdersConsumer)
	gob.Register(goesiv1.GetMarketsRegionIdOrders200Ok{})
}

func marketOrdersConsumer(s *Hammer, parameter interface{}) {
	regionID := parameter.(int32)
	var page int32 = 1

	for {
		orders, _, err := s.esi.V1.MarketApi.GetMarketsRegionIdOrders("all", regionID, map[string]interface{}{"page": page})
		if err != nil {
			log.Println(err)
			return
		} else if len(orders) == 0 { // end of the pages
			break
		}
		b, err := gobcoder.GobEncoder(orders)
		if err != nil {
			log.Println(err)
			return
		}

		err = s.nsq.Publish("marketOrders", b)
		if err != nil {
			log.Println(err)
			return
		}
	}

	return
}
