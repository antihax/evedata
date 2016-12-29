package views

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"

	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("wars", "GET", "/activeWars", activeWarsPage)
	evedata.AddRoute("wars", "GET", "/J/activeWars", activeWars)
}

func activeWarsPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	p := newPage(s, r, "Active Wars")

	templates.Templates = template.Must(template.ParseFiles("templates/wars.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func activeWars(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	v, err := models.GetActiveWarList()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}
