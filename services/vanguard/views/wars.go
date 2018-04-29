package views

import (
	"net/http"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("wars", "GET", "/activeWars",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"wars.html",
				time.Hour*24*31,
				newPage(r, "Active Wars"))
		})
	vanguard.AddRoute("wars", "GET", "/lostFighters",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"lostFighters.html",
				time.Hour*24*31,
				newPage(r, "Lost Fighters in HighSec"))
		})
	vanguard.AddRoute("wars", "GET", "/lossesInHighsec",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"lossesInHighSec.html",
				time.Hour*24*31,
				newPage(r, "Ship Losses in Highsec"))
		})

	vanguard.AddRoute("wars", "GET", "/J/activeWars", activeWars)
	vanguard.AddRoute("wars", "GET", "/J/lostFighters", lostFighters)
	vanguard.AddRoute("wars", "GET", "/J/lossesInHighsec", lossesInHighsec)
}

func activeWars(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetActiveWarList()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}

func lostFighters(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetCorporationAssetsInSpaceLostFightersHighsec()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func lossesInHighsec(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetLossesInHighsec()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour*24)
}
