package views

import (
	"net/http"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("GET", "/activeWars",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"wars.html",
				time.Hour*24*31,
				newPage(r, "Active Wars"))
		})

	vanguard.AddRoute("GET", "/J/activeWars", activeWars)
}

func activeWars(w http.ResponseWriter, r *http.Request) {
	v, err := models.GetActiveWarList()
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}
