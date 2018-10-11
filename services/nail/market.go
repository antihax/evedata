package nail

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/evedata/internal/datapackages"

	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("structureOrders", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddConcurrentHandlers(s.wait(nsq.HandlerFunc(s.structureOrderHandler)), 50)
	})

	AddHandler("marketHistory", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddConcurrentHandlers(s.wait(nsq.HandlerFunc(s.marketHistoryHandler)), 50)
	})
}

func (s *Nail) structureOrderHandler(message *nsq.Message) error {
	b := datapackages.StructureOrders{}
	err := gobcoder.GobDecoder(message.Body, &b)
	if err != nil {
		return err
	}

	if len(b.Orders) == 0 {
		return nil
	}

	values := []string{}
	count := 0
	for _, e := range b.Orders {
		count++
		var buy byte
		if e.IsBuyOrder {
			buy = 1
		} else {
			buy = 0
		}
		values = append(values, fmt.Sprintf("(%d,%f,%d,%d,%d,%d,%d,%q,%d,%d,evedata.regionIDByStructureID(%d),UTC_TIMESTAMP(),1)",
			e.OrderId, e.Price, e.VolumeRemain, e.TypeId, e.VolumeTotal, e.MinVolume,
			buy, e.Issued.UTC().Format("2006-01-02 15:04:05"), e.Duration, e.LocationId, e.LocationId))
		if count >= 80 {
			s.doMarketOrders(values)
			values = []string{}
			count = 0
		}
	}

	if err := s.doMarketOrders(values); err != nil {
		log.Println(err)
		return err
	}

	return s.inQueue.SetWorkExpire("evedata_entity", int64(allianceID), 43200)

}

func (s *Nail) doMarketOrders(values []string) error {
	if len(values) == 0 {
		return nil
	}
	stmt := fmt.Sprintf(`INSERT INTO evedata.market (orderID, price, remainingVolume, typeID, enteredVolume, minVolume, bid, issued, duration, stationID, regionID, reported, private)
	VALUES %s
	ON DUPLICATE KEY UPDATE price=VALUES(price), remainingVolume=VALUES(remainingVolume), issued=VALUES(issued), duration=VALUES(duration), reported=VALUES(reported);
		`, strings.Join(values, ",\n"))

	if err := s.doSQL(stmt); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Nail) marketHistoryHandler(message *nsq.Message) error {
	b := datapackages.MarketHistory{}
	err := gobcoder.GobDecoder(message.Body, &b)
	if err != nil {
		return err
	}
	var values []string
	ignoreBefore := time.Now().UTC().Add(-time.Hour * 24 * 60)

	for _, e := range b.History {
		orderDate, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			return err
		}

		if orderDate.After(ignoreBefore) {
			values = append(values, fmt.Sprintf("(%q,%f,%f,%f,%d,%d,%d,%d)",
				e.Date, e.Lowest, e.Highest, e.Average,
				e.Volume, e.OrderCount, b.TypeID, b.RegionID))
		}
	}

	if len(values) == 0 {
		return nil
	}

	stmt := fmt.Sprintf("INSERT INTO evedata.market_history (date, low, high, mean, quantity, orders, itemID, regionID) VALUES \n%s ON DUPLICATE KEY UPDATE date=date", strings.Join(values, ",\n"))

	return s.doSQL(stmt)
}
