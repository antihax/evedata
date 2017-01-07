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
	evedata.AddRoute("assets", "GET", "/assets", assetsPage)
	evedata.AddRoute("assets", "GET", "/U/assets", assetsAPI)
}

func assetsPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "Asset Information")
	templates.Templates = template.Must(template.ParseFiles("templates/assets.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func assetsAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 5*60)

	if s.Values["characterID"] == nil || s.Values["characterID"] == 0 {
		return http.StatusForbidden, nil
	}

	characterID := s.Values["characterID"].(int64)

	assets, err := models.GetAssets(characterID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(assets)

	return 200, nil
}
