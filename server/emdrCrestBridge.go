package evedata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmcvetta/napping"
)

var stations map[int64]int64

func goEMDRCrestBridge(c *AppContext) {

	type marketRegions struct {
		RegionID   int64  `db:"regionID"`
		RegionName string `db:"regionName"`
	}

	type marketTypes struct {
		TypeID   int64  `db:"typeID"`
		TypeName string `db:"typeName"`
	}

	regions := []marketRegion{}
	err := c.Db.Select(&regions, `
		SELECT 	regionID, regionName 
		FROM 	mapRegions 
		WHERE 	regionID < 11000000 
			AND regionID NOT IN(10000001, 10000017, 10000019, 10000004);
	`)
	if err != nil {
		log.Fatal("EMDRCrestBridge:", err)
	}
	log.Printf("EMDRCrestBridge: Loaded %d Regions", len(regions))

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

	stations = make(map[int64]int64)
	rows, err := c.Db.Query(`
		SELECT stationID, solarSystemID 
		FROM staStations;
	`)
	for rows.Next() {

		var stationID int64
		var systemID int64

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

	// Throttle Crest Requests
	rate := time.Second / 8
	throttle := time.Tick(rate)

	// semaphore to prevent runaways
	sem := make(chan bool, c.Conf.EMDRCrestBridge.MaxGoRoutines)

	// CREST Session
	crest := napping.Session{}

	for {
		// loop through all regions
		for _, r := range regions {
			// and each item per region
			for _, t := range types {
				<-throttle // impliment throttle

				sem <- true
				go func() {
					defer func() { <-sem }()
					// Process Market History
					h := marketHistory{}
					url := fmt.Sprintf("https://public-crest.eveonline.com/market/%d/types/%d/history/", r.RegionID, t.TypeID)
					response, err := crest.Get(url, nil, &h, nil)
					if err != nil {
						log.Printf("EMDRCrestBridge: %s %+v", err, response.RawText)
						return
					}
					if response.Status() == 200 {
						if len(h.Items) > 0 {
							go postHistory(h, c, t.TypeID, r.RegionID)
						}
					}
				}()

				sem <- true
				go func() {
					defer func() { <-sem }()
					// Process Market Buy Orders
					b := marketOrders{}
					url := fmt.Sprintf("https://public-crest.eveonline.com/market/%d/orders/buy/?type=https://public-crest.eveonline.com/types/%d/", r.RegionID, t.TypeID)
					response, err := crest.Get(url, nil, &b, nil)
					if err != nil {
						log.Printf("EMDRCrestBridge: %s %+v", err, response.RawText)
						return
					}
					if response.Status() == 200 {
						if len(b.Items) > 0 {
							go postOrders(b, c, 1, t.TypeID, r.RegionID)
						}
					}
				}()

				sem <- true
				go func() {
					defer func() { <-sem }()
					// Process Market Sell Orders
					s := marketOrders{}
					url := fmt.Sprintf("https://public-crest.eveonline.com/market/%d/orders/sell/?type=https://public-crest.eveonline.com/types/%d/", r.RegionID, t.TypeID)
					response, err := crest.Get(url, nil, &s, nil)
					if err != nil {
						log.Printf("EMDRCrestBridge: %s %+v", err, response.RawText)
						return
					}
					if response.Status() == 200 {
						if len(s.Items) > 0 {

							go postOrders(s, c, 0, t.TypeID, r.RegionID)
						}
					}
				}()
			}
		}
	}
}

func postHistory(h marketHistory, c *AppContext, typeID int64, regionID int64) {
	if c.Conf.EMDRCrestBridge.Import {
		historyUpdate, err := c.Db.Prepare(`
			INSERT IGNORE INTO market_history 
				(date, low, high, mean, quantity, orders, itemID, regionID) 
				VALUES(?,?,?,?,?,?,?,?);
		`)
		defer historyUpdate.Close()
		if err != nil {
			log.Fatalf("EMDRCrestBridge: %s", err)
		}

		tx, err := c.Db.Begin()
		for _, e := range h.Items {
			tx.Stmt(historyUpdate).Exec(e.Date, e.LowPrice, e.HighPrice, e.AvgPrice, typeID, regionID)
		}
		tx.Commit()
	}

	if c.Conf.EMDRCrestBridge.Upload {
		u := newUUDIFHeader()
		u.ResultType = "history"
		u.Columns = []string{"date", "orders", "quantity", "low", "high", "average"}

		u.Rowsets = make([]rowsetsUUDIF, 1)

		u.Rowsets[0].RegionID = regionID
		u.Rowsets[0].TypeID = typeID
		u.Rowsets[0].GeneratedAt = time.Now()

		u.Rowsets[0].Rows = make([][]interface{}, len(h.Items))

		for i, e := range h.Items {
			u.Rowsets[0].Rows[i] = make([]interface{}, 6)
			u.Rowsets[0].Rows[i][0] = e.Date + "+00:00"
			u.Rowsets[0].Rows[i][1] = e.OrderCount
			u.Rowsets[0].Rows[i][2] = e.Volume
			u.Rowsets[0].Rows[i][3] = e.LowPrice
			u.Rowsets[0].Rows[i][4] = e.HighPrice
			u.Rowsets[0].Rows[i][5] = e.AvgPrice
		}

		enc, err := json.Marshal(u)
		if err != nil {
			log.Println("EMDRCrestBridge:", err)
		}
		postUUDIF(c.Conf.EMDRCrestBridge.URL, enc)
	}
}

func postOrders(o marketOrders, c *AppContext, buy int, typeID int64, regionID int64) {

	if c.Conf.EMDRCrestBridge.Import {
		orderUpdate, err := c.Db.Prepare(`
					INSERT INTO market
						(orderID, price, remainingVolume, typeID, enteredVolume, minVolume, bid, issued, duration, stationID, regionID, systemID, reported)
						VALUES(?,?,?,?,?,?,?,?,?,?,?,?,NOW())
						ON DUPLICATE KEY UPDATE price=VALUES(price),
												remainingVolume=VALUES(remainingVolume),
												issued=VALUES(issued),
												duration=VALUES(duration),
												reported=VALUES(reported),
												done=0;
				`)
		defer orderUpdate.Close()

		orderMark, err := c.Db.Prepare(`
					UPDATE market SET done = 1 WHERE regionID = ? AND typeID =?
				`)
		defer orderMark.Close()

		if err != nil {
			log.Fatalf("EMDRCrestBridge: %s", err)
		}

		// Mark orders complete
		tx, err := c.Db.Begin()
		tx.Stmt(orderMark).Exec(regionID, typeID)

		// Add or update orders
		for _, e := range o.Items {
			tx.Stmt(orderUpdate).Exec(e.ID, e.Price, e.Volume, typeID, e.VolumeEntered, e.MinVolume, buy, e.Issued, e.Duration, e.Location.ID, regionID, stations[e.Location.ID])
		}

		err = tx.Commit()

		if err != nil {
			log.Println("EMDRCrestBridge:", err)
		}
	}

	if c.Conf.EMDRCrestBridge.Upload {
		u := newUUDIFHeader()
		u.ResultType = "orders"
		u.Columns = []string{"price", "volRemaining", "range", "orderID", "volEntered", "minVolume", "bid", "issueDate", "duration", "stationID", "solarSystemID"}

		u.Rowsets = make([]rowsetsUUDIF, 1)

		u.Rowsets[0].RegionID = regionID
		u.Rowsets[0].TypeID = typeID
		u.Rowsets[0].GeneratedAt = time.Now()

		u.Rowsets[0].Rows = make([][]interface{}, len(o.Items))

		for i, e := range o.Items {

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

			u.Rowsets[0].Rows[i] = make([]interface{}, 11)
			u.Rowsets[0].Rows[i][0] = e.Price
			u.Rowsets[0].Rows[i][1] = e.Volume
			u.Rowsets[0].Rows[i][2] = r
			u.Rowsets[0].Rows[i][3] = e.ID
			u.Rowsets[0].Rows[i][4] = e.VolumeEntered
			u.Rowsets[0].Rows[i][5] = e.MinVolume
			u.Rowsets[0].Rows[i][6] = e.Buy
			u.Rowsets[0].Rows[i][7] = e.Issued + "+00:00"
			u.Rowsets[0].Rows[i][8] = e.Duration
			u.Rowsets[0].Rows[i][9] = e.Location.ID
			u.Rowsets[0].Rows[i][10] = stations[e.Location.ID]
		}

		enc, err := json.Marshal(u)
		if err != nil {
			log.Println("EMDRCrestBridge:", err)
		}

		//log.Printf("%s", enc)
		postUUDIF(c.Conf.EMDRCrestBridge.URL, enc)
	}
}

func postUUDIF(url string, j []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("EMDRCrestBridge:", string(body))
		log.Println("EMDRCrestBridge:", string(resp.Status))

	}
}

func newUUDIFHeader() marketUUDIF {
	n := marketUUDIF{}

	n.Version = "0.1"

	n.Generator.Name = "EveData.Org"
	n.Generator.Version = "0.025a"

	n.UploadKeys = make([]uploadKeysUUDIF, 1)
	n.UploadKeys[0] = uploadKeysUUDIF{"EveData.Org", "TheCheeseIsBree"}

	n.CurrentTime = time.Now()

	return n
}

type rowsetsUUDIF struct {
	GeneratedAt time.Time       `json:"generatedAt"`
	RegionID    int64           `json:"regionID"`
	TypeID      int64           `json:"typeID"`
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

type marketHistory struct {
	TotalCount_Str string
	Items          []struct {
		OrderCount int64
		LowPrice   float64
		HighPrice  float64
		AvgPrice   float64
		Volume     int64
		Date       string
	}
	PageCount  int64
	TotalCount int64
}

type marketOrders struct {
	Items []struct {
		Buy           bool
		Issued        string
		Price         float64
		VolumeEntered int64
		MinVolume     int64
		Volume        int64
		Range         string
		Duration      int64
		ID            int64
		Location      struct {
			ID   int64
			Name string
		}
		Type struct {
			ID   int64
			Name string
		}
	}
	PageCount  int64
	TotalCount int64
}
