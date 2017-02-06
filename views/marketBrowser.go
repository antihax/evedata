package views

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("marketBrowser", "GET", "/", marketBrowser)
	evedata.AddRoute("searchItems", "GET", "/J/searchItems", searchitemsPage)
	evedata.AddRoute("marketSellRegionItems", "GET", "/J/marketSellRegionItems", MarketSellRegionItems)
	evedata.AddRoute("marketBuyRegionItems", "GET", "/J/marketBuyRegionItems", MarketBuyRegionItems)
}

// marketBrowser generates.... stuff
func marketBrowser(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	p := newPage(r, "EVE Online Market Browser")

	templates.Templates = template.Must(template.ParseFiles("templates/marketBrowser.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

type marketItemList struct {
	TypeID     int64  `db:"typeID"`
	TypeName   string `db:"typeName"`
	Categories string `db:"Categories"`
	Count      int64
}

// ARows bridge for old version
type ARows struct {
	Rows *[]marketItemList `json:"rows"`
}

func searchitemsPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var q string
	q = r.FormValue("q")

	if len(q) < 2 {
		return http.StatusInternalServerError, errors.New("Query too short")
	}

	mIL := []marketItemList{}

	// [BENCHMARK] 0.078 sec / 0.000 sec
	err := c.Db.Select(&mIL, `SELECT  T.typeID, typeName, CONCAT_WS(',', G5.marketGroupName, G4.marketGroupName, G3.marketGroupName, G2.marketGroupName, G.marketGroupName) AS Categories, count(*) AS count
           FROM invTypes T 
           LEFT JOIN invMarketGroups G on T.marketGroupID = G.marketGroupID
           LEFT JOIN invMarketGroups G2 on G.parentGroupID = G2.marketGroupID
           LEFT JOIN invMarketGroups G3 on G2.parentGroupID = G3.marketGroupID
           LEFT JOIN invMarketGroups G4 on G3.parentGroupID = G4.marketGroupID
           LEFT JOIN invMarketGroups G5 on G4.parentGroupID = G5.marketGroupID

           WHERE published=1 AND T.marketGroupID IS NOT NULL AND typeName LIKE ?
           GROUP BY T.typeID
           ORDER BY typeName
           LIMIT 100`, "%"+q+"%")

	if err != nil {
		return http.StatusInternalServerError, err
	}

	var mRows ARows

	mRows.Rows = &mIL

	encoder := json.NewEncoder(w)
	encoder.Encode(mRows)

	return 200, nil
}

/******************************************************************************
 * marketSellRegionItems JSON query
 *****************************************************************************/
type marketItems struct {
	StationName string `db:"stationName" json:"stationName"`
	StationID   string `db:"stationID"   json:"stationID"   `
	Quantity    string `db:"quantity"    json:"quantity"   `
	Price       string `db:"price"       json:"price"      `
}

// Rows is a list of rows for JSON conversion
type Rows struct {
	Rows *[]marketItems `json:"rows"`
}

const (
	highSec = 1 << iota
	lowSec  = 1 << iota
	nullSec = 1 << iota
)

// MarketRegionItems Query market orders for a user specified
func marketRegionItems(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, buy bool) (int, error) {
	var (
		mRows         Rows
		err           error
		secFilter     string
		secFilterPass int
	)

	mR := []marketItems{}

	regionID, err := strconv.Atoi(r.FormValue("regionID"))
	if err != nil {
		regionID = 0
	}

	itemID, err := strconv.Atoi(r.FormValue("itemID"))
	if err != nil {
		return 500, err
	}

	secFlags, err := strconv.Atoi(r.FormValue("secflags"))
	if err != nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(mR)
		return 200, nil
	}

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
		sql := `SELECT  remainingVolume AS quantity, price, stationName, M.stationID
        	                    FROM    evedata.market M
                             	INNER JOIN staStations S ON S.stationID=M.stationID
                             	INNER JOIN mapSolarSystems Sy ON Sy.solarSystemID = S.solarSystemID
                             	WHERE      bid=? AND
                                      	   typeID = ? AND (` + secFilter + `)`
		err = c.Db.Select(&mR, sql, buy, itemID)
	} else {
		err = c.Db.Select(&mR, `SELECT  remainingVolume AS quantity, price, stationName, M.stationID
        	                    FROM    evedata.market M
                             	INNER JOIN staStations S ON S.stationID=M.stationID
								INNER JOIN mapSolarSystems Sy ON Sy.solarSystemID = S.solarSystemID
                             	WHERE      bid=? AND
                                      	   M.regionID = ? AND
                                      	   typeID = ? AND (`+secFilter+`)`, buy, regionID, itemID, secFilter)
	}

	if err != nil {
		return 500, err
	}

	mRows.Rows = &mR

	encoder := json.NewEncoder(w)
	encoder.Encode(mR)
	return 200, nil
}

// MarketSellRegionItems Query market sell orders for a user specified
// regionID and itemID query string and return JSON to the user
func MarketSellRegionItems(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	return marketRegionItems(c, w, r, false)
}

// MarketBuyRegionItems Query market buy orders for a user specified
// regionID and itemID query string and return JSON to the user
func MarketBuyRegionItems(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	return marketRegionItems(c, w, r, true)
}
