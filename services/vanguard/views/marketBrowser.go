package views

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("marketBrowser", "GET", "/marketBrowser",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"marketBrowser.html",
				time.Hour*24*31,
				newPage(r, "Market Browser"))
		})
	vanguard.AddRoute("searchMarketItems", "GET", "/J/searchMarketItems", searchMarketItemsAPI)
	vanguard.AddRoute("marketSellRegionItems", "GET", "/J/marketSellRegionItems", MarketSellRegionItems)
	vanguard.AddRoute("marketBuyRegionItems", "GET", "/J/marketBuyRegionItems", MarketBuyRegionItems)
}

func searchMarketItemsAPI(w http.ResponseWriter, r *http.Request) {
	var q string
	q = r.FormValue("q")

	if len(q) < 2 {
		httpErr(w, errors.New("Query too short"))
		return
	}

	v, err := models.SearchMarketNames(q)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
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

	v, err := models.MarketRegionItems(regionID, itemID, secFlags, buy)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
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
