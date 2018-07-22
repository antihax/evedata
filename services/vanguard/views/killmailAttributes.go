package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("GET", "/killmailAttributes",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w, "killmailAttributes.html", time.Hour*24*31, newPage(r, "Killmail Attribute Browser"))
		})

	vanguard.AddRoute("GET", "/J/killmailAttributesAPI", killmailAttributesAPI)
	vanguard.AddRoute("GET", "/J/offensiveGroups", offensiveGroupsAPI)
}

func offensiveGroupsAPI(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetArbitrageCalculatorStations()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}

func killmailAttributesAPI(w http.ResponseWriter, r *http.Request) {
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

	renderJSON(w, v, time.Hour)
}
