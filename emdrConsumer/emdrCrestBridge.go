package emdrConsumer

import (
	"bytes"
	"evedata/appContext"
	"evedata/esi"
	"evedata/eveapi"
	"evedata/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var stations map[int64]int64

// Run the bridge between CREST and Eve Market Data Relay.
// Optionally import to the database
func GoEMDRCrestBridge(c *appContext.AppContext) {
	type regionKey struct {
		RegionID int64
		TypeID   int64
	}

	type marketRegions struct {
		RegionID   int64  `db:"regionID"`
		RegionName string `db:"regionName"`
	}

	type marketTypes struct {
		TypeID   int64  `db:"typeID"`
		TypeName string `db:"typeName"`
	}

	type marketOrders struct {
		regionID int64
		orders   *[]esi.GetMarketsRegionIdOrders200Ok
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
	postChannel := make(chan []byte, 20)
	historyChannel := make(chan *eveapi.MarketTypeHistoryCollectionV1, 20)
	orderChannel := make(chan marketOrders, 20)

	if c.Conf.EMDRCrestBridge.Upload {
		for i := 0; i < 10; i++ {
			// Don't spawn them all at once.
			time.Sleep(time.Second / 2)
			go func() {
				for {
					msg := <-postChannel

					response, err := http.Post(c.Conf.EMDRCrestBridge.URL, "application/json", bytes.NewBuffer(msg))
					if err != nil {
						log.Println("EMDRCrestBridge:", err)
					} else {
						if response.Status != "200 OK" {
							log.Println("EMDRCrestBridge:", string(response.Status))
						}
						// Must read everything to close the body and reuse connection
						ioutil.ReadAll(response.Body)
						response.Body.Close()
					}
				}
			}()
		}
	}

	if c.Conf.EMDRCrestBridge.Import {
		go func() {
			for {
				h := <-historyChannel
				if len(h.Items) == 0 {
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

					for _, e := range h.Items {
						values = append(values, fmt.Sprintf("('%s',%f,%f,%f,%d,%d,%d,%d)",
							e.Date, e.LowPrice, e.HighPrice, e.AvgPrice,
							e.Volume, e.OrderCount, h.TypeID, h.RegionID))
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
		}()
		go func() {
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
		}()
	}

	// limit concurrent requests as to not hog the available connections.
	// Eventually the buffers will become the limiting factors.
	limiter := make(chan bool, 20)
	for {
		// loop through all regions
		for _, r := range regions {
			// and each item per region
			for _, t := range types {
				limiter <- true
				go func(l chan bool) {
					defer func(l chan bool) { <-l }(l)
					// Process Market Buy Orders
					b, err := c.ESI.MarketApi.GetMarketsRegionIdOrders(r.RegionID, "all", t.TypeID, nil, nil)
					order := marketOrders{r.RegionID, &b}
					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
						return
					}

					if c.Conf.EMDRCrestBridge.Import {
						orderChannel <- order
					}
				}(limiter)

				limiter <- true
				go func(l chan bool) {
					defer func(l chan bool) { <-l }(l)
					rk := regionKey{r.RegionID, t.TypeID}

					// Process Market History
					h, err := c.EVE.MarketTypeHistoryV1ByID(rk.RegionID, rk.TypeID)
					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
						return
					}

					if c.Conf.EMDRCrestBridge.Import {
						historyChannel <- h
					}
				}(limiter)
			}
		}
	}
}
