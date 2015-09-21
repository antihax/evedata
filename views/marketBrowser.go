package views

import (
	"encoding/json"
	"errors"
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"
	"strconv"
)

func init() {
	evedata.AddRoute(evedata.Route{"marketBrowser", "GET", "/", marketBrowser})
	evedata.AddRoute(evedata.Route{"searchItems", "GET", "/J/searchItems", searchitemsPage})
	evedata.AddRoute(evedata.Route{"marketSellRegionItems", "GET", "/J/marketSellRegionItems", MarketSellRegionItems})
	evedata.AddRoute(evedata.Route{"marketBuyRegionItems", "GET", "/J/marketBuyRegionItems", MarketBuyRegionItems})
}

// marketBrowser generates.... stuff
func marketBrowser(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	p := Page{
		Title: "EVE Online Market Browser",
	}

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

func searchitemsPage(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	var q string
	q = r.FormValue("q")

	if len(q) < 2 {
		return http.StatusInternalServerError, errors.New("Query too short")
	}

	mIL := []marketItemList{}

	err := c.Db.Select(&mIL, `SELECT  T.typeID, typeName, CONCAT_WS(',', G5.marketGroupName, G4.marketGroupName, G3.marketGroupName, G2.marketGroupName, G.marketGroupName) AS Categories, count(*) AS count
           FROM market M
           INNER JOIN invTypes T ON M.typeID = T.typeID
           LEFT JOIN invMarketGroups G on T.marketGroupID = G.marketGroupID
           LEFT JOIN invMarketGroups G2 on G.parentGroupID = G2.marketGroupID
           LEFT JOIN invMarketGroups G3 on G2.parentGroupID = G3.marketGroupID
           LEFT JOIN invMarketGroups G4 on G3.parentGroupID = G4.marketGroupID
           LEFT JOIN invMarketGroups G5 on G4.parentGroupID = G5.marketGroupID

           WHERE done=0 AND T.marketGroupID IS NOT NULL AND typeName LIKE ?
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
	HighSec = 1 << iota
	LowSec  = 1 << iota
	NullSec = 1 << iota
)

// MarketRegionItems Query market orders for a user specified
func marketRegionItems(c *evedata.AppContext, w http.ResponseWriter, r *http.Request, buy bool) (int, error) {
	var (
		regionID      int
		itemID        int
		secFlags      int
		mRows         Rows
		err           error
		secFilter     string
		secFilterPass int
	)

	mR := []marketItems{}

	regionID, err = strconv.Atoi(r.FormValue("regionID"))
	if err != nil {
		regionID = 0
	}

	itemID, err = strconv.Atoi(r.FormValue("itemID"))
	if err != nil {
		return 500, err
	}

	secFlags, err = strconv.Atoi(r.FormValue("secflags"))
	if err != nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(mR)
		return 200, nil
	}

	if secFlags&HighSec != 0 {
		secFilterPass++
		secFilter += "round(Sy.security,1) >= 0.5"
	}
	if secFlags&LowSec != 0 {
		secFilterPass++
		if secFilterPass > 1 {
			secFilter += " OR "
		}
		secFilter += "round(Sy.security,1) BETWEEN 0.1 AND 0.4"
	}
	if secFlags&NullSec != 0 {
		secFilterPass++
		if secFilterPass > 1 {
			secFilter += " OR "
		}

		secFilter += "round(Sy.security,1) <= 0 "
	}

	if regionID == 0 {
		sql := `SELECT  format(remainingVolume, 0) AS quantity, format(price, 2) as price, stationName, M.stationID
        	                    FROM    market M
                             	INNER JOIN staStations S ON S.stationID=M.stationID
                             	INNER JOIN mapSolarSystems Sy ON Sy.solarSystemID = M.systemID
                             	WHERE      done=0 AND
                                	       bid=? AND
                                      	   typeID = ? AND (` + secFilter + `) ORDER BY price ASC`
		err = c.Db.Select(&mR, sql, buy, itemID)
	} else {
		err = c.Db.Select(&mR, `SELECT  format(remainingVolume, 0) AS quantity, format(price, 2) as price, stationName, M.stationID
        	                    FROM    market M
                             	INNER JOIN staStations S ON S.stationID=M.stationID
                             	INNER JOIN mapSolarSystems Sy ON Sy.solarSystemID = M.systemID
                             	WHERE      done=0 AND
                                	       bid=? AND
                                      	   M.regionID = ? AND
                                      	   typeID = ? AND (`+secFilter+`) ORDER BY price ASC`, buy, regionID, itemID, secFilter)
	}

	if err != nil {
		return 500, err
	}

	mRows.Rows = &mR

	encoder := json.NewEncoder(w)
	encoder.Encode(mR)
	return 200, nil
}

func MarketSellRegionItems(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	return marketRegionItems(c, w, r, false)
}

// MarketBuyRegionItems Query market buy orders for a user specified
// regionID and itemID query string and return JSON to the user
func MarketBuyRegionItems(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	return marketRegionItems(c, w, r, true)
}
