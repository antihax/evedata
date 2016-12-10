package emdrConsumer

import (
	"evedata/appContext"
	"evedata/esi"
	"evedata/models"
	"fmt"
	"log"
	"strings"
)

var stations map[int64]int64

type marketOrders struct {
	regionID int32
	orders   *[]esi.GetMarketsRegionIdOrders200Ok
}

type marketHistory struct {
	regionID int32
	typeID   int32
	history  *[]esi.GetMarketsRegionIdHistory200Ok
}

// Run the bridge between CREST and Eve Market Data Relay.
// Optionally import to the database
func GoEMDRCrestBridge(c *appContext.AppContext) {
	type regionKey struct {
		RegionID int32
		TypeID   int32
	}

	type marketRegions struct {
		RegionID   int32  `db:"regionID"`
		RegionName string `db:"regionName"`
	}

	type marketTypes struct {
		TypeID   int32  `db:"typeID"`
		TypeName string `db:"typeName"`
	}

	// Obtain a list of regions which have market stations
	regions, err := models.GetMarketRegions()
	if err != nil {
		log.Fatal("EMDRCrestBridge:", err)
	}
	log.Printf("EMDRCrestBridge: Loaded %d Regions", len(regions))

	// Obtain list of types available on the market
	types := []marketTypes{}
	err = c.Db.Select(&types, `
		SELECT 	typeID, typeName 
		FROM 	invTypes 
		WHERE 	marketGroupID IS NOT NULL 
			AND typeID < 250000;
	`)
	if err != nil {
		log.Fatal("EMDRCrestBridge:", err)
	}
	log.Printf("EMDRCrestBridge: Loaded %d items", len(types))

	// Get a list of stations
	stations = make(map[int64]int64)
	rows, err := c.Db.Query(`
		SELECT stationID, solarSystemID 
		FROM staStations;
	`)
	for rows.Next() {
		var stationID, systemID int64
		if err := rows.Scan(&stationID, &systemID); err != nil {
			log.Fatal("EMDRCrestBridge: ", err)
		}
		stations[stationID] = systemID
	}
	rows.Close()

	if err != nil {
		log.Fatal("EMDRCrestBridge: ", err)
	}
	log.Printf("EMDRCrestBridge: Loaded %d stations", len(stations))

	// Build buffers for posting to the database and
	historyChannel := make(chan marketHistory, 50)
	orderChannel := make(chan marketOrders, 50)

	// Start the consumers.
	go historyConsumer(historyChannel, c)
	go orderConsumer(orderChannel, c)

	// limit concurrent requests as to not hog the available connections.
	// Eventually the buffers will become the limiting factors.
	limiter := make(chan bool, 5)
	for {

		// Update the market data
		log.Printf("EMDRCrestBridge: updateMarket.")
		_, err := c.Db.Exec(`call updateMarket;`)
		if err != nil {
			log.Printf("EMDRCrestBridge: Failed updateMarket: %v", err)
			return
		}

		// loop through all regions
		for _, r := range regions {
			limiter <- true
			go func(l chan bool) {
				defer func(l chan bool) { <-l }(l)

				// Start at page 1
				var page int32 = 1
				// Process Market Buy Orders
				for {
					b, _, err := c.ESI.MarketApi.GetMarketsRegionIdOrders(r.RegionID, "all", nil, page, nil)
					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
						return
					} else if len(b) == 0 { // end of the pages
						break
					}

					// Post the orders
					order := marketOrders{r.RegionID, &b}
					orderChannel <- order

					// Next page
					page++
				}
			}(limiter)
			// and each item per region
			for _, t := range types {

				limiter <- true
				go func(l chan bool) {
					defer func(l chan bool) { <-l }(l)

					// Process Market History
					h, _, err := c.ESI.MarketApi.GetMarketsRegionIdHistory(r.RegionID, t.TypeID, nil)

					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
						return
					}

					hist := marketHistory{r.RegionID, t.TypeID, &h}
					historyChannel <- hist
				}(limiter)
			}
		}
	}
}

func orderConsumer(orderChannel chan marketOrders, c *appContext.AppContext) {
	{
		for {
			o := <-orderChannel
			// Add or update orders
			if len(*o.orders) == 0 {
				continue
			}

			for {
				tx, err := c.Db.Begin()
				if err != nil {
					log.Printf("EMDRCrestBridge: %s", err)
					continue
				}

				var values []string
				for _, e := range *o.orders {
					var buy byte
					if e.IsBuyOrder == true {
						buy = 1
					} else {
						buy = 0
					}
					values = append(values, fmt.Sprintf("(%d,%f,%d,%d,%d,%d,%d,'%s',%d,%d,%d,%d,UTC_TIMESTAMP())",
						e.OrderId, e.Price, e.VolumeRemain, e.TypeId, e.VolumeTotal, e.MinVolume,
						buy, e.Issued.UTC().Format("2006-01-02 15:04:05"), e.Duration, e.LocationId, o.regionID, stations[e.LocationId]))
				}

				stmt := fmt.Sprintf(`INSERT INTO market (orderID, price, remainingVolume, typeID, enteredVolume, minVolume, bid, issued, duration, stationID, regionID, systemID, reported)
						VALUES %s
						ON DUPLICATE KEY UPDATE price=VALUES(price),
							remainingVolume=VALUES(remainingVolume),
							issued=VALUES(issued),
							duration=VALUES(duration),
							reported=VALUES(reported),
							done=0;
							`, strings.Join(values, ",\n"))

				_, err = tx.Exec(stmt)

				if err != nil {
					log.Printf("EMDRCrestBridge: %s", err)
					tx.Rollback()
					break
				}

				err = tx.Commit()
				if err != nil {
					log.Printf("EMDRCrestBridge: %s", err)
					continue
				}
				break // success
			}
		}
	}
}
func historyConsumer(historyChannel chan marketHistory, c *appContext.AppContext) {
	for {
		h := <-historyChannel
		if len(*h.history) == 0 {
			continue
		}

		// Loop until the transaction passes
		for {
			tx, err := c.Db.Begin()
			if err != nil {
				log.Printf("EMDRCrestBridge: %s", err)
				break
			}
			var values []string

			for _, e := range *h.history {
				values = append(values, fmt.Sprintf("('%s',%f,%f,%f,%d,%d,%d,%d)",
					e.Date, e.Lowest, e.Highest, e.Average,
					e.Volume, e.OrderCount, h.typeID, h.regionID))
			}

			stmt := fmt.Sprintf("INSERT IGNORE INTO market_history (date, low, high, mean, quantity, orders, itemID, regionID) VALUES \n %s", strings.Join(values, ",\n"))

			_, err = tx.Exec(stmt)
			if err != nil {
				tx.Rollback()
				log.Printf("EMDRCrestBridge: %s", err)
				break
			}

			err = tx.Commit()
			if err != nil {
				log.Printf("EMDRCrestBridge: %s", err)
				break
			}
			break // success
		}
	}
}
