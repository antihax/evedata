package evedata

import (
	"fmt"
	"log"
	"time"

	"github.com/jmcvetta/napping"
)

func goEMDRCrestBridge(c *AppContext) {
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

	rate := time.Second / 15
	throttle := time.Tick(rate)
	for {
		for _, r := range regions {
			for _, t := range types {
				<-throttle
				res := history{}
				url := fmt.Sprintf("https://public-crest.eveonline.com/market/%d/types/%d/history/", r.RegionID, t.TypeID)
				response, err := napping.Get(url, nil, &res, nil)
				if err != nil {
					panic(err)
				}
				if response.Status() == 200 {
					log.Printf("%+v\n", res)
				}
			}
		}
	}

}

type marketRegions struct {
	RegionID   int64  `db:"regionID"`
	RegionName string `db:"regionName"`
}

type marketTypes struct {
	TypeID   int64  `db:"typeID"`
	TypeName string `db:"typeName"`
}

type history struct {
	TotalCount_Str string
	Items          []struct {
		Volume_str     string
		OrderCount     int64
		LowPrice       float64
		HighPrice      float64
		AvgPrice       float64
		Volume         int64
		OrderCount_str string
		Date           string
	}
	PageCount     int64
	PageCount_str string
	TotalCount    int64
}
