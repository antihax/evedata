package views

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/evedata"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("marketBrowser", "GET", "/marketBrowser", marketBrowser)
	evedata.AddRoute("searchMarketItems", "GET", "/J/searchMarketItems", searchMarketItemsAPI)
	evedata.AddRoute("marketSellRegionItems", "GET", "/J/marketSellRegionItems", MarketSellRegionItems)
	evedata.AddRoute("marketBuyRegionItems", "GET", "/J/marketBuyRegionItems", MarketBuyRegionItems)
}

// marketBrowser generates.... stuff
func marketBrowser(w http.ResponseWriter, r *http.Request) {
	p := newPage(r, "Market Browser")

	templates.Templates = template.Must(template.ParseFiles("templates/marketBrowser.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)
	if err != nil {
		httpErr(w, err)
		return
	}
}

// ARows bridge for old version
type ARows struct {
	Rows *[]models.MarketItemList `json:"rows"`
}

func searchMarketItemsAPI(w http.ResponseWriter, r *http.Request) {

	var q string
	q = r.FormValue("q")

	if len(q) < 2 {
		httpErr(w, errors.New("Query too short"))
		return
	}

	mIL, err := models.SearchMarketNames(q)
	if err != nil {
		httpErr(w, err)
		return
	}

	var mRows ARows
	mRows.Rows = &mIL

	json.NewEncoder(w).Encode(mRows)
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
func marketRegionItems(w http.ResponseWriter, r *http.Request, buy bool) {
	var (
		err           error
		secFilter     string
		secFilterPass int
	)

	c := evedata.GlobalsFromContext(r.Context())

	mR := []marketItems{}

	regionID, err := strconv.Atoi(r.FormValue("regionID"))
	if err != nil {
		regionID = 0
	}

	itemID, err := strconv.Atoi(r.FormValue("itemID"))
	if err != nil {
		httpErr(w, err)
		return
	}

	secFlags, err := strconv.Atoi(r.FormValue("secflags"))
	if err != nil {
		httpErr(w, err)
		return
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
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(mR)
}

// MarketSellRegionItems Query market sell orders for a user specified
// regionID and itemID query string and return JSON to the user
func MarketSellRegionItems(w http.ResponseWriter, r *http.Request) {
	marketRegionItems(w, r, false)
}

// MarketBuyRegionItems Query market buy orders for a user specified
// regionID and itemID query string and return JSON to the user
func MarketBuyRegionItems(w http.ResponseWriter, r *http.Request) {
	marketRegionItems(w, r, true)
}
