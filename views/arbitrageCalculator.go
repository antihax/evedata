package views

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("arbitrageCalculator", "GET", "/arbitrageCalculator", arbitrageCalculatorPage)
	evedata.AddRoute("arbitrageCalculatorStations", "GET", "/J/arbitrageCalculatorStations", arbitrageCalculatorStations)
	evedata.AddRoute("arbitrageCalculator", "GET", "/J/arbitrageCalculator", arbitrageCalculator)
}

func arbitrageCalculatorPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60*24)
	p := newPage(r, "Arbitrage Calculator")

	templates.Templates = template.Must(template.ParseFiles("templates/arbitrageCalculator.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		httpErr(w, err)
		return
	}
}

func arbitrageCalculatorStations(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*30)
	v, err := models.GetArbitrageCalculatorStations()
	if err != nil {
		httpErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}

func arbitrageCalculator(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*30)

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

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}
