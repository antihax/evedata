package evedata

import (
	"bytes"
	"encoding/json"
	"evedata/appContext"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/jmcvetta/napping"
)

var stations map[int64]int64

// Temporary Hack
func getKills(r int64, client napping.Session, c *appContext.AppContext) {
	type kills struct {
		KillID        int
		SolarSystemID int
		KillTime      string
		MoonID        int
		Victim        struct {
			CharacterID   int
			CorporationID int
			AllianceID    int
		}
		ZKB struct {
			Hash string
		}
	}

	url := fmt.Sprintf("https://zkillboard.com/api/kills/regionID/%d/", r)
	h := []kills{}
	response, err := client.Get(url, nil, &h, nil)
	if err != nil {
		log.Printf("EMDRCrestBridge: %s", err)
		return
	}
	if response.Status() == 200 {
		tx, err := c.Db.Begin()
		if err != nil {
			log.Printf("EMDRCrestBridge: %s", err)
			return
		}
		for _, e := range h {
			tx.Stmt(c.Bridge.KillInsert).Exec(e.KillID, e.SolarSystemID, e.KillTime, e.MoonID, e.Victim.CharacterID, e.Victim.CorporationID, e.Victim.AllianceID, e.ZKB.Hash)
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("EMDRCrestBridge: %s", err)
			return
		}
	}
}

func goEMDRCrestBridge(c *appContext.AppContext) {

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

	// FanOut response channel for posters
	postChannel := make(chan []byte)

	client := c.HTTPClient

	if c.Conf.EMDRCrestBridge.Import {
		var err error
		c.Bridge.HistoryUpdate, err = c.Db.Prepare(`
			INSERT IGNORE INTO market_history 
				(date, low, high, mean, quantity, orders, itemID, regionID) 
				VALUES(?,?,?,?,?,?,?,?);
		`)
		if err != nil {
			log.Fatalf("EMDRCrestBridge: %s", err)
		}

		c.Bridge.KillInsert, err = c.Db.Prepare(`
			INSERT IGNORE INTO killmails (id,solarSystemID,killTime,moonID,victimCharacterID,victimCorporationID,victimAllianceID,hash)
			VALUES(?,?,?,?,?,?,?,?);
		`)
		if err != nil {
			log.Fatalf("EMDRCrestBridge: %s", err)
		}

		c.Bridge.OrderUpdate, err = c.Db.Prepare(`
					INSERT INTO market
						(orderID, price, remainingVolume, typeID, enteredVolume, minVolume, bid, issued, duration, stationID, regionID, systemID, reported)
						VALUES(?,?,?,?,?,?,?,?,?,?,?,?,UTC_TIMESTAMP())
						ON DUPLICATE KEY UPDATE price=VALUES(price),
												remainingVolume=VALUES(remainingVolume),
												issued=VALUES(issued),
												duration=VALUES(duration),
												reported=VALUES(reported),
												done=0;
				`)
		if err != nil {
			log.Fatalf("EMDRCrestBridge: %s", err)
		}

		c.Bridge.OrderMark, err = c.Db.Prepare(`
					UPDATE market SET done = 1 WHERE regionID = ? AND typeID =?
				`)

		if err != nil {
			log.Fatalf("EMDRCrestBridge: %s", err)
		}
	}

	go func() {
		for i := 0; i < 11; i++ {
			// Don't spawn them all at once.
			time.Sleep(time.Second / 2)
			go func() {
				for {
					msg := <-postChannel

					response, err := client.Post(c.Conf.EMDRCrestBridge.URL, "application/json", bytes.NewBuffer(msg))
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
	}()

	// Throttle Crest Requests
	rate := time.Second / 30
	throttle := time.Tick(rate)

	// semaphore to prevent runaways
	sem := make(chan bool, c.Conf.EMDRCrestBridge.MaxGoRoutines)
	sem2 := make(chan bool, c.Conf.EMDRCrestBridge.MaxGoRoutines)

	// CREST Session
	crest := napping.Session{Client: client}

	for {
		// loop through all regions
		for _, r := range regions {
			<-throttle // impliment throttle
			sem2 <- true
			go func() {
				defer func() { <-sem2 }()
				// Process Market Buy Orders
				b := marketOrders{}
				url := fmt.Sprintf("https://crest-tq.eveonline.com/market/%d/orders/all/", r.RegionID)

				go getKills(r.RegionID, crest, c)

				response, err := crest.Get(url, nil, &b, nil)
				if err != nil {
					log.Printf("EMDRCrestBridge: %s", err)
					return
				}
				if response.Status() == 200 {
					sem <- true
					go postOrders(sem, postChannel, b, c, r.RegionID)

					next := b.Next.Href
					for len(next) > 0 {
						n := marketOrders{}
						res, err := crest.Get(next, nil, &n, nil)
						if err != nil {
							log.Printf("EMDRCrestBridge: %s", err)
							return
						}
						if res.Status() == 200 {
							sem <- true
							go postOrders(sem, postChannel, n, c, r.RegionID)
						}
						next = n.Next.Href
					}
				}
			}()

			// and each item per region
			for _, t := range types {
				<-throttle // impliment throttle
				sem2 <- true

				rk := regionKey{r.RegionID, t.TypeID}
				go func() {
					defer func() { <-sem2 }()
					// Process Market History
					h := marketHistory{}
					url := fmt.Sprintf("https://crest-tq.eveonline.com/market/%d/history/?type=https://crest-tq.eveonline.com/inventory/types/%d/", rk.RegionID, rk.TypeID)

					response, err := crest.Get(url, nil, &h, nil)
					if err != nil {
						log.Printf("EMDRCrestBridge: %s", err)
						return
					}
					if response.Status() == 200 {
						sem <- true
						go postHistory(sem, postChannel, h, c, rk.RegionID, rk.TypeID)
					}
				}()
			}
		}
	}
}

func postHistory(sem chan bool, postChan chan []byte, h marketHistory, c *appContext.AppContext, regionID int64, typeID int64) {
	defer func() { <-sem }()
	if c.Conf.EMDRCrestBridge.Import {

		tx, err := c.Db.Begin()
		if err != nil {
			log.Printf("EMDRCrestBridge: %s", err)
			return
		}
		for _, e := range h.Items {
			//(date, low, high, mean, quantity, orders, itemID, regionID)
			//fmt.Printf("%s %f %f %f %d %d %d %d\n", e.Date, e.LowPrice, e.HighPrice, e.AvgPrice, e.Volume, e.OrderCount, typeID, regionID)
			_, err := tx.Stmt(c.Bridge.HistoryUpdate).Exec(e.Date, e.LowPrice, e.HighPrice, e.AvgPrice, e.Volume, e.OrderCount, typeID, regionID)
			if err != nil {
				log.Printf("EMDRCrestBridge: %s", err)
				return
			}
		}
		err = tx.Commit()
		if err != nil {
			log.Println("EMDRCrestBridge:", err)
		}
	}

	if c.Conf.EMDRCrestBridge.Upload {
		u := newUUDIFHeader()
		u.ResultType = "history"
		u.Columns = []string{"date", "orders", "quantity", "low", "high", "average"}

		rowsets := make(map[int64]rowsetsUUDIF)

		edit := rowsets[typeID]
		edit.RegionID = regionID
		edit.TypeID = typeID
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
		rowsets[typeID] = edit
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
}

func postOrders(sem chan bool, postChan chan []byte, o marketOrders, c *appContext.AppContext, regionID int64) {
	defer func() { <-sem }()
	if c.Conf.EMDRCrestBridge.Import {

		// Mark orders complete
		tx, err := c.Db.Begin()
		if err != nil {
			log.Printf("EMDRCrestBridge: %s", err)
			return
		}
		// Add or update orders
		first := false
		for _, e := range o.Items {
			if first {
				_, err = tx.Stmt(c.Bridge.OrderMark).Exec(regionID, e.Type)
				if err != nil {
					log.Printf("EMDRCrestBridge: %s", err)
					return
				}
			}
			_, err = tx.Stmt(c.Bridge.OrderUpdate).Exec(e.ID, e.Price, e.Volume, e.Type, e.VolumeEntered, e.MinVolume, e.Buy, e.Issued, e.Duration, e.StationID, regionID, stations[e.StationID])
			if err != nil {
				log.Printf("EMDRCrestBridge: %s", err)
				return
			}
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
			edit.RegionID = regionID
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
			row[7] = e.Issued + "+00:00"
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
}

func newUUDIFHeader() marketUUDIF {
	n := marketUUDIF{}

	n.Version = "0.1"

	n.Generator.Name = "EveData.Org"
	n.Generator.Version = "0.1a"

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
	Next       struct {
		Href string
	}
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
		Type          int64
		StationID     int64
	}
	PageCount  int64
	TotalCount int64
	Next       struct {
		Href string
	}
}
