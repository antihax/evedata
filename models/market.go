package models

import (
	"fmt"
	"time"
)

type MarketHistory struct {
	Date     time.Time `db:"date" json:"date"`
	Low      float64   `db:"low" json:"low"`
	High     float64   `db:"high" json:"high"`
	Open     float64   `db:"open" json:"open"`
	Close    float64   `db:"close" json:"close"`
	Quantity int64     `db:"quantity" json:"quantity"`
}

// [BENCHMARK] 0.407 sec / 0.421 sec [TODO] Optimize
func GetMarketHistory(itemID int64, regionID int32) ([]MarketHistory, error) {
	s := []MarketHistory{}
	if err := database.Select(&s, `
		SELECT H.date, H.low, H.high, H.mean AS close, Y.mean AS open, H.quantity 
		FROM evedata.market_history H
		INNER JOIN evedata.market_history Y ON H.date = DATE_SUB(Y.date, INTERVAL 1 DAY) 
			AND H.regionID = Y.regionID 
			AND H.itemID = Y.itemID
		WHERE H.regionID = ? AND H.itemID = ? AND H.quantity > 10
	`, regionID, itemID); err != nil {
		return nil, err
	}
	return s, nil
}

type ArbitrageCalculatorStations struct {
	StationName string `db:"stationName" json:"stationName" `
	StationID   string `db:"stationID" json:"stationID" `
}

// [BENCHMARK] 0.015 sec / 0.000 sec
func GetArbitrageCalculatorStations() ([]ArbitrageCalculatorStations, error) {
	s := []ArbitrageCalculatorStations{}
	if err := database.Select(&s, `
		SELECT stationID, stationName
			FROM    evedata.marketStations
			WHERE 	Count > 4000
			ORDER BY stationName
	`); err != nil {
		return nil, err
	}
	return s, nil
}

type ArbitrageCalculator struct {
	TypeID   int64  `db:"typeID" json:"typeID" `
	TypeName string `db:"typeName" json:"typeName" `
	Volume   int64  `db:"volume" json:"volume" `
	Buys     int64  `db:"buys" json:"buys" `
	Margin   string `db:"margin" json:"margin" `
}

func GetArbitrageCalculator(stationID int64, minVolume int64, maxPrice int64, brokersFee float64, tax float64, method string) ([]ArbitrageCalculator, error) {
	type buys struct {
		TypeID   int64   `db:"typeID"`
		TypeName string  `db:"typeName"`
		Buys     int64   `db:"buys"`
		Volume   int64   `db:"volume"`
		Price    float64 `db:"price"`
	}
	type sells struct {
		TypeID int64   `db:"typeID"`
		Price  float64 `db:"price"`
	}

	errc := make(chan error)

	b := []buys{}

	go func() {
		// [BENCHMARK] 0.432 sec / 0.016 sec
		err := database.Select(&b, `
		SELECT  market.typeID AS typeID, typeName, count(*) as buys, ROUND(market_vol.quantity / 2) as volume, ROUND(max(price) + (max(price) * ?),2) AS price
		FROM    evedata.market, invTypes, evedata.market_vol
		WHERE   market.typeID = market_vol.itemID AND
		        market.regionID = market_vol.regionID AND
		        market.typeID = invTypes.typeID AND
		        market.stationID = ?
		        AND bid = 1
		        AND market_vol.quantity /2 > ?
		GROUP BY market.typeID
		HAVING price < ?`, brokersFee, stationID, minVolume, maxPrice)

		errc <- err
	}()

	sellOrders := make(map[int64]sells)
	s := []sells{}
	go func() {
		// [BENCHMARK] 0.297 sec / 0.000 sec
		err := database.Select(&s, `
		SELECT  typeID, ROUND(min(price) - (min(price) * ?) - (min(price) * ?),2) AS price
		FROM    evedata.market
		WHERE   market.stationID = ?
		        AND bid = 0
		        AND price < ?
		        GROUP BY typeID`, brokersFee, tax, stationID, maxPrice)

		// Add broker fee to all orders.
		for _, order := range s {
			sellOrders[order.TypeID] = order
		}
		errc <- err
	}()

	margins := []ArbitrageCalculator{}

	// clear the error channel
	for i := 0; i < 2; i++ {
		err := <-errc
		if err != nil {
			return nil, err
		}
	}

	for _, buyOrder := range b {
		sellOrder, ok := sellOrders[buyOrder.TypeID]

		if !ok {
			continue
		}
		if method == "delta" {
			if (sellOrder.Price - buyOrder.Price) > 500000 {
				newOrder := ArbitrageCalculator{}

				newOrder.Buys = buyOrder.Buys
				newOrder.Volume = ((int64)(sellOrder.Price+buyOrder.Price) / 2) * buyOrder.Volume
				newOrder.TypeID = buyOrder.TypeID
				newOrder.TypeName = buyOrder.TypeName
				newOrder.Margin = fmt.Sprintf("%.2f", sellOrder.Price-buyOrder.Price)

				margins = append(margins, newOrder)
			}
		} else if method == "percentage" {
			if (sellOrder.Price / buyOrder.Price) > 1.05 {
				newOrder := ArbitrageCalculator{}

				newOrder.Buys = buyOrder.Buys
				newOrder.Volume = ((int64)(sellOrder.Price+buyOrder.Price) / 2) * buyOrder.Volume
				newOrder.TypeID = buyOrder.TypeID
				newOrder.TypeName = buyOrder.TypeName
				newOrder.Margin = fmt.Sprintf("%.2f", sellOrder.Price/buyOrder.Price)

				margins = append(margins, newOrder)
			}
		}
	}

	return margins, nil
}

type MarketRegion struct {
	RegionID   int32  `db:"regionID"`
	RegionName string `db:"regionName"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
// Anywhere can now have a public market.
func GetMarketRegions() ([]MarketRegion, error) {
	v := []MarketRegion{}
	err := database.Select(&v, `
		SELECT 	regionID, regionName 
		FROM 	mapRegions 
		WHERE regionID < 11000000;
	`)
	if err != nil {
		return nil, err
	}
	return v, nil
}

type MarketType struct {
	TypeID   int32  `db:"typeID"`
	TypeName string `db:"typeName"`
}

// [BENCHMARK] 0.000 sec / 0.047 sec
func GetMarketTypes() ([]MarketType, error) {
	v := []MarketType{}
	err := database.Select(&v, `
		SELECT 	typeID, typeName 
		FROM 	invTypes 
		WHERE 	marketGroupID IS NOT NULL
	`)
	if err != nil {
		return nil, err
	}
	return v, nil
}

type MarketItems struct {
	StationName string `db:"stationName" json:"stationName"`
	StationID   string `db:"stationID"   json:"stationID"   `
	Quantity    string `db:"quantity"    json:"quantity"   `
	Price       string `db:"price"       json:"price"      `
}

const (
	highSec = 1 << iota
	lowSec  = 1 << iota
	nullSec = 1 << iota
)

func MarketRegionItems(regionID int, itemID int, secFlags int, buy bool) ([]MarketItems, error) {
	var (
		secFilter     string
		secFilterPass int
		err           error
	)

	mR := []MarketItems{}

	if secFlags&highSec != 0 {
		secFilterPass++
		secFilter += "round(Sy.security,1) >= 0.5"
	}

	if secFlags&lowSec != 0 {
		secFilterPass++
		if secFilterPass > 1 {
			secFilter += " OR "
		}
		secFilter += "round(Sy.security,1) BETWEEN 0.1 AND 0.4"
	}

	if secFlags&nullSec != 0 {
		secFilterPass++
		if secFilterPass > 1 {
			secFilter += " OR "
		}

		secFilter += "round(Sy.security,1) <= 0 "
	}

	if regionID == 0 {
		err = database.Select(&mR, `
			SELECT  remainingVolume AS quantity, price, stationName, M.stationID
				FROM    evedata.market M
				INNER JOIN staStations S ON S.stationID=M.stationID
				INNER JOIN mapSolarSystems Sy ON Sy.solarSystemID = S.solarSystemID
				WHERE	bid=? AND
						typeID = ? AND (`+secFilter+`);
			`, buy, itemID)
	} else {
		err = database.Select(&mR, `
			SELECT  remainingVolume AS quantity, price, stationName, M.stationID
				FROM    evedata.market M
				INNER JOIN staStations S ON S.stationID=M.stationID
				INNER JOIN mapSolarSystems Sy ON Sy.solarSystemID = S.solarSystemID
				WHERE	bid=? AND
						M.regionID = ? AND
						typeID = ? AND (`+secFilter+`);
			`, buy, regionID, itemID)
	}

	return mR, err
}
