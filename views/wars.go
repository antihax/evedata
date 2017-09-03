package views

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/antihax/evedata/evedata"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("wars", "GET", "/activeWars", activeWarsPage)
	evedata.AddRoute("wars", "GET", "/J/activeWars", activeWars)
	evedata.AddRoute("wars", "GET", "/lostFighters", lostFightersPage)
	evedata.AddRoute("wars", "GET", "/J/lostFighters", lostFighters)
	evedata.AddRoute("wars", "GET", "/lossesInHighsec", lossesInHighsecPage)
	evedata.AddRoute("wars", "GET", "/J/lossesInHighsec", lossesInHighsec)
}

func activeWarsPage(w http.ResponseWriter, r *http.Request) {
	p := newPage(r, "Active Wars")

	templates.Templates = template.Must(template.ParseFiles("templates/wars.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		httpErr(w, err)
		return
	}
}

func activeWars(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	v, err := models.GetActiveWarList()
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}

func lostFightersPage(w http.ResponseWriter, r *http.Request) {
	p := newPage(r, "Fighters Lost in HighSec")

	templates.Templates = template.Must(template.ParseFiles("templates/lostFighters.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		httpErr(w, err)
		return
	}
}

func lostFighters(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	v, err := models.GetCorporationAssetsInSpaceLostFightersHighsec()
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}

func lossesInHighsecPage(w http.ResponseWriter, r *http.Request) {
	p := newPage(r, "Fighters Lost in HighSec")

	templates.Templates = template.Must(template.ParseFiles("templates/lossesInHighSec.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		httpErr(w, err)
		return
	}
}

func lossesInHighsec(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60*24)
	v, err := models.GetLossesInHighsec()
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}
