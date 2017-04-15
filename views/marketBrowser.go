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

	json.NewEncoder(w).Encode(mIL)
}

func marketRegionItems(w http.ResponseWriter, r *http.Request, buy bool) {
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

	mR, err := models.MarketRegionItems(regionID, itemID, secFlags, buy)
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
