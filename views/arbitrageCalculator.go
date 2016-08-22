package views

import (
	"encoding/json"
	"evedata/appContext"
	"evedata/models"
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("arbitrageCalculator", "GET", "/arbitrageCalculator", arbitrageCalculatorPage)
	evedata.AddRoute("arbitrageCalculatorStations", "GET", "/J/arbitrageCalculatorStations", arbitrageCalculatorStations)
	evedata.AddRoute("arbitrageCalculator", "GET", "/J/arbitrageCalculator", arbitrageCalculator)
}

func arbitrageCalculatorPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	p := newPage(s, r, "Arbitrage Calculator")

	templates.Templates = template.Must(template.ParseFiles("templates/arbitrageCalculator.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func arbitrageCalculatorStations(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	v, err := models.GetArbitrageCalculatorStations()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}

func arbitrageCalculator(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	hours, err := strconv.ParseInt(r.FormValue("hours"), 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

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

	brokersFee = brokersFee / 100
	tax = tax / 100

	v, err := models.GetArbitrageCalculator(hours, stationID, minVolume, maxPrice, brokersFee, tax)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}
