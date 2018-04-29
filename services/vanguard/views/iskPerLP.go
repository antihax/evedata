package views

import (
	"net/http"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("iskPerLP", "GET", "/iskPerLP",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"iskPerLP.html",
				time.Hour*24*31,
				newPage(r, "ISK Per Loyalty Point"))
		})
	vanguard.AddRoute("iskPerLP", "GET", "/iskPerLPByConversion",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"iskPerLPByConversion.html",
				time.Hour*24*31,
				newPage(r, "ISK Per Loyalty Point - All Entities"))
		})
	vanguard.AddRoute("iskPerLPCorps", "GET", "/J/iskPerLPCorps", iskPerLPCorps)
	vanguard.AddRoute("iskPerLP", "GET", "/J/iskPerLP", iskPerLP)
	vanguard.AddRoute("iskPerLP", "GET", "/J/iskPerLPByConversion", iskPerLPByConversion)
}

func iskPerLPCorps(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetISKPerLPCorporations()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}

func iskPerLP(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("corp")
	v, err := models.GetISKPerLP(q)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}

func iskPerLPByConversion(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetISKPerLPByConversion()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}
