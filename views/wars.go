package views

import (
	"encoding/json"
	"evedata/appContext"
	"evedata/models"
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"

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
	v, err := models.GetActiveWarList()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}
