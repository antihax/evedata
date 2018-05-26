package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("bubblePlacer", "GET", "/bubblePlacer",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"bubblePlacer.html",
				time.Hour*24*31,
				newPage(r, "Bubble-O-Matic 9002"))
		})
	vanguard.AddRoute("iskPerLP", "GET", "/iskPerLPByConversion",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"iskPerLPByConversion.html",
				time.Hour*24*31,
				newPage(r, "ISK Per Loyalty Point - All Entities"))
		})
	vanguard.AddRoute("nullSystems", "GET", "/J/nullSystems", nullSystems)
	vanguard.AddRoute("systemCelestials", "GET", "/J/systemCelestials", systemCelestials)

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
