package emdrConsumer

import (
	"bytes"
	"encoding/json"
	"evedata/appContext"
	"evedata/eveapi"
	"evedata/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	orderChannel := make(chan *eveapi.MarketOrderCollectionSlimV1, 20)

	if c.Conf.EMDRCrestBridge.Upload {
		for i := 0; i < 4; i++ {
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
				for _, e := range h.Items {
					_, err := c.Db.Exec(`INSERT IGNORE INTO market_history 
						(date, low, high, mean, quantity, orders, itemID, regionID) 
						VALUES(?,?,?,?,?,?,?,?)`,
						e.Date, e.LowPrice, e.HighPrice, e.AvgPrice,
						e.Volume, e.OrderCount, h.TypeID, h.RegionID)
					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
					}
				}
			}
		}()
		go func() {
			for {
				o := <-orderChannel
				// Add or update orders
				first := false
				for _, e := range o.Items {
					if first {
						_, err := c.Db.Exec(`UPDATE market SET done = 1 WHERE regionID = ? AND typeID =? LIMIT 50`, o.RegionID, e.Type)
						if err != nil {
							log.Printf("EMDRCrestBridge: %s", err)
						}
					}
					_, err := c.Db.Exec(`INSERT INTO market
						(orderID, price, remainingVolume, typeID, enteredVolume, minVolume, bid, issued, duration, stationID, regionID, systemID, reported)
						VALUES(?,?,?,?,?,?,?,?,?,?,?,?,UTC_TIMESTAMP())
						ON DUPLICATE KEY UPDATE price=VALUES(price),
												remainingVolume=VALUES(remainingVolume),
												issued=VALUES(issued),
												duration=VALUES(duration),
												reported=VALUES(reported),
												done=0;`,
						e.ID, e.Price, e.Volume, e.Type, e.VolumeEntered, e.MinVolume,
						e.Buy, e.Issued.UTC(), e.Duration, e.StationID, o.RegionID, stations[e.StationID])
					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
					}
				}
			}
		}()
	}

	// limit concurrent requests as to not hog the available connections.
	// Eventually the buffers will become the limiting factors.
	limiter := make(chan bool, 4)
	for {
		// loop through all regions
		for _, r := range regions {
			limiter <- true
			go func(l chan bool) {
				defer func(l chan bool) { <-l }(l)
				// Process Market Buy Orders
				b, err := c.EVE.MarketOrdersSlim(r.RegionID, 1)
				if err != nil {
					log.Printf("EMDRCrestBridge: %s", err)
					return
				}

				for ; b != nil; b, err = b.NextPage() {
					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
						return
					}
					if c.Conf.EMDRCrestBridge.Upload {
						postOrders(postChannel, b, c)
					}
					if c.Conf.EMDRCrestBridge.Import {
						orderChannel <- b
					}
				}
			}(limiter)

			// and each item per region
			for _, t := range types {
				limiter <- true
				go func(l chan bool) {
					defer func(l chan bool) { <-l }(l)
					rk := regionKey{r.RegionID, t.TypeID}

					// Process Market History
					h, err := c.EVE.MarketTypeHistoryByID(rk.RegionID, rk.TypeID)

					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
						return
					}

					if c.Conf.EMDRCrestBridge.Upload {
						postHistory(postChannel, h, c)
					}
					if c.Conf.EMDRCrestBridge.Import {
						historyChannel <- h
					}
				}(limiter)
			}
		}
	}
}

func postHistory(postChan chan []byte, h *eveapi.MarketTypeHistoryCollectionV1, c *appContext.AppContext) {
	u := newUUDIFHeader()
	u.ResultType = "history"
	u.Columns = []string{"date", "orders", "quantity", "low", "high", "average"}

	rowsets := make(map[int64]rowsetsUUDIF)

	edit := rowsets[h.TypeID]
	edit.RegionID = h.RegionID
	edit.TypeID = h.TypeID
	edit.GeneratedAt = time.Now()

	edit.Rows = make([][]interface{}, len(h.Items))

	for i, e := range h.Items {
		edit.Rows[i] = make([]interface{}, 6)
		edit.Rows[i][0] = e.Date + "+00:00"
		edit.Rows[i][1] = e.OrderCount
		edit.Rows[i][2] = e.Volume
		edit.Rows[i][3] = e.LowPrice
		edit.Rows[i][4] = e.HighPrice
		edit.Rows[i][5] = e.AvgPrice
	}
	rowsets[h.TypeID] = edit
	for _, v := range rowsets {
		u.Rowsets = append(u.Rowsets, v)
	}

	enc, err := json.Marshal(u)

	if err != nil {
		log.Println("EMDRCrestBridge:", err)
	} else {
		postChan <- enc
	}
}

func postOrders(postChan chan []byte, o *eveapi.MarketOrderCollectionSlimV1, c *appContext.AppContext) {
	u := newUUDIFHeader()
	u.ResultType = "orders"
	u.Columns = []string{"price", "volRemaining", "range", "orderID", "volEntered", "minVolume", "bid", "issueDate", "duration", "stationID", "solarSystemID"}
	rowsets := make(map[int64]rowsetsUUDIF)

	for _, e := range o.Items {
		var r int

		switch {
		case e.Range == "station":
			r = -1
		case e.Range == "solarsystem":
			r = 0
		case e.Range == "region":
			r = 32767
		default:
			r, _ = strconv.Atoi(e.Range)
		}

		edit := rowsets[e.Type]
		edit.RegionID = o.RegionID
		edit.GeneratedAt = time.Now()
		edit.TypeID = e.Type
		row := make([]interface{}, 11)
		row[0] = e.Price
		row[1] = e.Volume
		row[2] = r
		row[3] = e.ID
		row[4] = e.VolumeEntered
		row[5] = e.MinVolume
		row[6] = e.Buy
		row[7] = e.Issued.UTC()
		row[8] = e.Duration
		row[9] = e.StationID
		row[10] = stations[e.StationID]
		edit.Count++
		edit.Rows = append(edit.Rows, row)
		rowsets[e.Type] = edit

	}

	for _, v := range rowsets {
		u.Rowsets = append(u.Rowsets, v)
	}

	enc, err := json.Marshal(u)
	if err != nil {
		log.Println("EMDRCrestBridge:", err)
	} else {
		postChan <- enc
	}
}

func newUUDIFHeader() marketUUDIF {
	n := marketUUDIF{}

	n.Version = "0.2"

	n.Generator.Name = "EveData.Org"
	n.Generator.Version = "0.2q.343"

	n.UploadKeys = make([]uploadKeysUUDIF, 1)
	n.UploadKeys[0] = uploadKeysUUDIF{"EveData.Org", "TheCheeseIsBree"}

	n.CurrentTime = time.Now()

	return n
}

type rowsetsUUDIF struct {
	GeneratedAt time.Time       `json:"generatedAt"`
	RegionID    int64           `json:"regionID"`
	TypeID      int64           `json:"typeID"`
	Count       int64           `json:"-"`
	Rows        [][]interface{} `json:"rows"`
}

type uploadKeysUUDIF struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type marketUUDIF struct {
	ResultType string            `json:"resultType"`
	Version    string            `json:"version"`
	UploadKeys []uploadKeysUUDIF `json:"uploadKeys"`
	Generator  struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"generator"`
	Columns     []string       `json:"columns"`
	CurrentTime time.Time      `json:"currentTime"`
	Rowsets     []rowsetsUUDIF `json:"rowsets"`
}
