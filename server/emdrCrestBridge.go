package evedata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jmcvetta/napping"
)

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
			AND regionID NOT IN(10000017, 10000019, 10000004);
	`)
	if err != nil {
		log.Fatal("EMDRCrestBridge: %s", err)
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
		log.Fatal("EMDRCrestBridge: %s", err)
	}
	log.Printf("EMDRCrestBridge: Loaded %d items", len(types))

	// Throttle Crest Requests to 30 RPS
	rate := time.Second / 10
	throttle := time.Tick(rate)

	// semaphore to prevent runaways
	sem := make(chan bool, c.Conf.EMDRCrestBridge.MaxGoRoutines)

	for {
		// loop through all regions
		for _, r := range regions {
			// and each item per region
			for _, t := range types {
				<-throttle // impliment throttle

				h := marketHistory{}
				url := fmt.Sprintf("https://public-crest.eveonline.com/market/%d/types/%d/history/", r.RegionID, t.TypeID)
				response, err := napping.Get(url, nil, &h, nil)
				if err != nil {
					log.Printf("EMDRCrestBridge: %s", err)
				}
				if response.Status() == 200 {
					if len(h.Items) > 0 {
						sem <- true
						go postHistory(sem, h, c, t.TypeID, r.RegionID)
					}
				}
			}
		}
	}
}

// https://public-crest.eveonline.com/market/10000002/orders/buy/?type=https://public-crest.eveonline.com/types/683/

func postHistory(sem chan bool, h marketHistory, c *AppContext, typeID int64, regionID int64) {
	defer func() { <-sem }()

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
	n.Generator.Version = "The Cheese You Smell Is Gorgonzola"

	n.UploadKeys = make([]uploadKeysUUDIF, 1)
	n.UploadKeys[0] = uploadKeysUUDIF{"EDO", "SomethingMadeUp"}

	n.CurrentTime = time.Now()

	return n
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

type marketOrder struct {
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
