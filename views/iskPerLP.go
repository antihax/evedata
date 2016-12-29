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
	evedata.AddRoute("iskPerLP", "GET", "/iskPerLP", iskPerLPPage)
	evedata.AddRoute("iskPerLPCorpss", "GET", "/J/iskPerLPCorps", iskPerLPCorps)
	evedata.AddRoute("iskPerLP", "GET", "/J/iskPerLP", iskPerLP)
}

func iskPerLPPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "ISK Per Loyalty Point")

	templates.Templates = template.Must(template.ParseFiles("templates/iskPerLP.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func iskPerLPCorps(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	v, err := models.GetISKPerLPCorporations()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}

func iskPerLP(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*30)
	q := r.FormValue("corp")
	v, err := models.GetISKPerLP(q)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}
