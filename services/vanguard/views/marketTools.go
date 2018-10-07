package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/goesi"
)

func init() {
	vanguard.AddRoute("GET", "/marketUndervalue",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"marketUnderValue.html",
				time.Hour*24*31,
				newPage(r, "EVE Online Undervalued Market Items"))
		})
	vanguard.AddRoute("GET", "/marketStationStocker",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"marketStationStocker.html",
				time.Hour*24*31,
				newPage(r, "EVE Online Station Stocker"))
		})
	vanguard.AddRoute("GET", "/J/marketRegions", marketRegionsAPI)
	vanguard.AddRoute("GET", "/J/marketUndervalue", marketUnderValueAPI)
	vanguard.AddAuthRoute("GET", "/J/marketStationStocker", marketStationStockerAPI)
}

func marketRegionsAPI(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetMarketRegions()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour*24)
}

func marketUnderValueAPI(w http.ResponseWriter, r *http.Request) {
	marketRegionID, err := strconv.ParseInt(r.FormValue("marketRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	sourceRegionID, err := strconv.ParseInt(r.FormValue("sourceRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	destinationRegionID, err := strconv.ParseInt(r.FormValue("destinationRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	discount, err := strconv.ParseFloat(r.FormValue("discount"), 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	discount = discount / 100

	v, err := models.MarketUnderValued(marketRegionID, sourceRegionID, destinationRegionID, discount)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Minute)
}

func marketStationStockerAPI(w http.ResponseWriter, r *http.Request) {
	characterID := int32(0)
	s := vanguard.SessionFromContext(r.Context())
	ch, ok := s.Values["character"].(goesi.VerifyResponse)
	if ok {
		characterID = ch.CharacterID
	}

	marketRegionID, err := strconv.ParseInt(r.FormValue("marketRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	destinationRegionID, err := strconv.ParseInt(r.FormValue("destinationRegionID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	markup, err := strconv.ParseFloat(r.FormValue("markup"), 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	markup = markup / 100

	v, err := models.MarketStationStocker(characterID, marketRegionID, destinationRegionID, markup)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Minute)
}
