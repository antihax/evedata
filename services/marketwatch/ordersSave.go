package marketwatch

import (
	"log"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/antihax/goesi/esi"
)

/*
func (s *MarketWatch) orderChangesPump() {
	for {
		c := <-s.orderChangeChan

		order := squirrel.Insert("evedata.marketOrderHistory").Columns("orderID", "changed", "locationID", "typeID", "volumeChange", "volumeRemain", "price", "duration", "isBuyOrder")

		count := 0
		for _, g := range c {
			if g.VolumeChange > 0 && // Skip price changes
				g.Issued.Add(time.Hour*24*time.Duration(g.Duration)).
					After(time.Now().UTC()) { // Skip expired
				count++
				order = order.Values(g.OrderID, g.TimeChanged, g.LocationId, g.TypeID, g.VolumeChange, g.VolumeRemain, g.Price, g.Duration, g.IsBuyOrder)
			}
		}

		// Early out if there is nothing to change
		if count == 0 {
			continue
		}

		sqlq, args, err := order.ToSql()
		if err != nil {
			log.Println(err)
		} else {
			err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE orderID = orderID", args...)
			if err != nil {
				log.Println(err)
			}
		}
	}
}*/

func (s *MarketWatch) saveAdditions(c []esi.GetMarketsRegionIdOrders200Ok) {
	order := squirrel.Insert("evedata.market").Columns(
		"orderID", "regionID", "stationID", "typeID", "bid", "price", "minVolume",
		"remainingVolume", "enteredVolume", "issued", "duration", "reported",
	)
	start := time.Now()
	count := 0
	length := len(c)

	wg := sync.WaitGroup{}
	sem := make(chan bool, 8)

	for _, g := range c {
		count++
		order = order.Values(
			g.OrderId,
			squirrel.Expr("evedata.regionIDByStructureID(?)", g.LocationId),
			g.LocationId, g.TypeId, g.IsBuyOrder, g.Price, g.MinVolume,
			g.VolumeRemain, g.VolumeTotal, g.Issued, g.Duration, squirrel.Expr("UTC_TIMESTAMP()"),
		)

		// If we have 500 or reached the end, dump what we have
		if count%500 == 0 || count == length {
			sqlq, args, err := order.ToSql()
			if err != nil {
				log.Println(err)
			} else {
				wg.Add(1)
				sem <- true
				go func(sqlq string, args []interface{}, start time.Time, count int) {
					defer func() { <-sem; wg.Done() }()
					err = s.doSQL(sqlq+`ON DUPLICATE KEY UPDATE price=VALUES(price),
				remainingVolume=VALUES(remainingVolume),
				issued=VALUES(issued), duration=VALUES(duration),
				reported=VALUES(reported);`, args...)
					if err != nil {
						log.Println(err)
					}
				}(sqlq, args, start, count)
			}
			// restart the query
			start = time.Now()
			order = squirrel.Insert("evedata.market").Columns(
				"orderID", "regionID", "stationID", "typeID", "bid", "price", "minVolume",
				"remainingVolume", "enteredVolume", "issued", "duration", "reported",
			)
		}
	}
	wg.Wait()
}

func (s *MarketWatch) saveDeletions(c []OrderChange) {
	or := squirrel.Or{}
	for _, g := range c {
		or = append(or, squirrel.Eq{"orderID": g.OrderID})
	}
	order := squirrel.Delete("evedata.market").Where(or)

	sqlq, args, err := order.ToSql()
	if err != nil {
		log.Println(err)
	} else {
		err = s.doSQL(sqlq, args...)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *MarketWatch) saveChanges(c []OrderChange) {
	for _, g := range c {
		sqlq, args, err := squirrel.Update("evedata.market").
			Set("price", g.Price).
			Set("issued", g.Issued).
			Set("remainingVolume", g.VolumeRemain).
			Set("duration", g.Duration).
			Set("reported", time.Now().UTC()).
			Where(squirrel.Eq{"orderID": g.OrderID}).ToSql()
		if err != nil {
			log.Println(err)
		} else {
			err = s.doSQL(sqlq, args...)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *MarketWatch) sqlPumps() {
	// Deal with change records in a separate routine as we don't have to worry about
	// collisions.
	//go s.orderChangesPump()

	for {
		select {
		case v := <-s.orderChan:
			s.saveAdditions(v)
		case v := <-s.orderDeleteChan:
			s.saveDeletions(v)
		case v := <-s.orderChangeChan:
			s.saveChanges(v)

		case v := <-s.contractChan:
			s.saveContractAdditions(v)
		case v := <-s.contractDeleteChan:
			s.saveContractChanges(v)
		case v := <-s.contractChangeChan:
			s.saveContractChanges(v)
		default:
			time.Sleep(time.Second)
		}
	}
}
