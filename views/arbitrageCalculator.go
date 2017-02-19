package views

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("arbitrageCalculator", "GET", "/arbitrageCalculator", arbitrageCalculatorPage)
	evedata.AddRoute("arbitrageCalculatorStations", "GET", "/J/arbitrageCalculatorStations", arbitrageCalculatorStations)
	evedata.AddRoute("arbitrageCalculator", "GET", "/J/arbitrageCalculator", arbitrageCalculator)
}

func arbitrageCalculatorPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60*24)
	p := newPage(r, "Arbitrage Calculator")

	templates.Templates = template.Must(template.ParseFiles("templates/arbitrageCalculator.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func arbitrageCalculatorStations(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*30)
	v, err := models.GetArbitrageCalculatorStations()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}

func arbitrageCalculator(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*30)

	stationID, err := strconv.ParseInt(r.FormValue("stationID"), 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	minVolume, err := strconv.ParseInt(r.FormValue("minVolume"), 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	maxPrice, err := strconv.ParseInt(r.FormValue("maxPrice"), 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	brokersFee, err := strconv.ParseFloat(r.FormValue("brokersFee"), 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tax, err := strconv.ParseFloat(r.FormValue("tax"), 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	method := r.FormValue("method")

	brokersFee = brokersFee / 100
	tax = tax / 100

	v, err := models.GetArbitrageCalculator(stationID, minVolume, maxPrice, brokersFee, tax, method)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}
