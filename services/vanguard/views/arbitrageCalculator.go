package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("GET", "/arbitrageCalculator",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "arbitrageCalculator.html", time.Hour*24*31, newPage(r, "Arbitrage Calculator"))
		})

	vanguard.AddRoute("GET", "/J/arbitrageCalculatorStations", arbitrageCalculatorStations)
	vanguard.AddRoute("GET", "/J/arbitrageCalculator", arbitrageCalculator)
}

func arbitrageCalculatorStations(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetArbitrageCalculatorStations()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Minute)
}

func arbitrageCalculator(w http.ResponseWriter, r *http.Request) {
	stationID, err := strconv.ParseInt(r.FormValue("stationID"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	minVolume, err := strconv.ParseInt(r.FormValue("minVolume"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	maxPrice, err := strconv.ParseInt(r.FormValue("maxPrice"), 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	brokersFee, err := strconv.ParseFloat(r.FormValue("brokersFee"), 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	tax, err := strconv.ParseFloat(r.FormValue("tax"), 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	method := r.FormValue("method")

	brokersFee = brokersFee / 100
	tax = tax / 100

	v, err := models.GetArbitrageCalculator(stationID, minVolume, maxPrice, brokersFee, tax, method)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Minute)
}
