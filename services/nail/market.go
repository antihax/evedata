package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/datapackages"

	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("marketOrders", spawnMarketConsumer)
}

func spawnMarketConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.marketHandler)))
}

func (s *Nail) marketHandler(message *nsq.Message) error {
	b := datapackages.MarketOrders{}
	err := gobcoder.GobDecoder(message.Body, &b)
	if err != nil {
		return err
	}
	log.Printf("%+v\n", b)
	if len(b.Orders) == 0 {
		return nil
	}

	var values []string
	for _, e := range b.Orders {
		var buy byte
		if e.IsBuyOrder == true {
			buy = 1
		} else {
			buy = 0
		}
		values = append(values, fmt.Sprintf("(%d,%f,%d,%d,%d,%d,%d,%q,%d,%d,%d,UTC_TIMESTAMP())",
			e.OrderId, e.Price, e.VolumeRemain, e.TypeId, e.VolumeTotal, e.MinVolume,
			buy, e.Issued.UTC().Format("2006-01-02 15:04:05"), e.Duration, e.LocationId, (int32)(b.RegionID)))
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.market (orderID, price, remainingVolume, typeID, enteredVolume, minVolume, bid, issued, duration, stationID, regionID, reported)
			VALUES %s
			ON DUPLICATE KEY UPDATE price=VALUES(price),
				remainingVolume=VALUES(remainingVolume),
				issued=VALUES(issued),
				duration=VALUES(duration),
				reported=VALUES(reported);
				`, strings.Join(values, ",\n"))
	log.Println(stmt)

	return s.DoSQL(stmt)
}
