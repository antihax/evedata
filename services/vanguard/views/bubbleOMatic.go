package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("GET", "/bubblePlacer",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"bubblePlacer.html",
				time.Hour*24*31,
				newPage(r, "Bubble-O-Matic 9002"))
		})
	vanguard.AddRoute("GET", "/battleFinder",
		func(w http.ResponseWriter, r *http.Request) {
			renderBlank(w,
				"battleFinder.html",
				time.Hour*24*31,
				newPage(r, "battlefinder"))
		})

	vanguard.AddRoute("GET", "/J/nullSystems", nullSystems)
	vanguard.AddRoute("GET", "/J/systemCelestials", systemCelestials)

}

func nullSystems(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetNullSystems()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour*24*31)
}

func systemCelestials(w http.ResponseWriter, r *http.Request) {

	solarSystemID, err := strconv.Atoi(r.FormValue("solarSystemID"))
	if err != nil {
		httpErr(w, err)
		return
	}

	v, err := models.GetSystemCelestials(int32(solarSystemID))
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour*24*31)
}
