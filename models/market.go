package models

import "fmt"

type ArbitrageCalculatorStations struct {
	StationName string `db:"stationName" json:"stationName" `
	StationID   string `db:"stationID" json:"stationID" `
}

func GetArbitrageCalculatorStations() ([]ArbitrageCalculatorStations, error) {
	s := []ArbitrageCalculatorStations{}
	if err := database.Select(&s, `
		SELECT stationID, stationName
			FROM    marketStations
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

func GetArbitrageCalculator(hours int64, stationID int64, minVolume int64, maxPrice int64, brokersFee float64, tax float64) ([]ArbitrageCalculator, error) {
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
		err := database.Select(&b, `
		SELECT  market.typeID AS typeID, typeName, count(*) as buys, ROUND(market_vol.quantity / 2) as volume, ROUND(max(price) + (max(price) * ?),2) AS price
		FROM    market, invTypes, market_vol
		WHERE   market.done = 0 AND
		        market.typeID = market_vol.itemID AND
		        market.regionID = market_vol.regionID AND
		        market.typeID = invTypes.typeID AND
		        reported >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL ? DAY_HOUR) AND
		        market.stationID = ?
		        AND bid = 1
		        AND done = 0 AND
		        market_vol.quantity /2 > ?
		GROUP BY market.typeID
		HAVING price < ?`, brokersFee, hours, stationID, minVolume, maxPrice)

		errc <- err
	}()

	sellOrders := make(map[int64]sells)
	s := []sells{}
	go func() {
		err := database.Select(&s, `
		SELECT  typeID, ROUND(min(price) - (min(price) * ?) - (min(price) * ?),2) AS price
		FROM    market
		WHERE   reported >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL ? DAY_HOUR) AND
		        market.stationID = ?
		        AND bid = 0
		        GROUP BY typeID`, brokersFee, tax, hours, stationID)

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

		if (sellOrder.Price - buyOrder.Price) > 500000 {
			newOrder := ArbitrageCalculator{}

			newOrder.Buys = buyOrder.Buys
			newOrder.Volume = buyOrder.Volume
			newOrder.TypeID = buyOrder.TypeID
			newOrder.TypeName = buyOrder.TypeName
			newOrder.Margin = fmt.Sprintf("%.2f", sellOrder.Price-buyOrder.Price)

			margins = append(margins, newOrder)
		}
	}

	return margins, nil
}

type MarketRegion struct {
	RegionID   int64  `db:"regionID"`
	RegionName string `db:"regionName"`
}

func GetMarketRegions() ([]MarketRegion, error) {
	v := []MarketRegion{}
	err := database.Select(&v, `
		SELECT 	regionID, regionName 
		FROM 	mapRegions 
		WHERE 	regionID < 11000000 
			AND regionID NOT IN(10000001, 10000017, 10000019, 10000004);
	`)
	if err != nil {
		return nil, err
	}
	return v, nil
}
