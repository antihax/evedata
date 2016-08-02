package views

import (
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute(evedata.Route{"ecmjam", "GET", "/ecmjam", ecmjamPage})
}

func ecmjamPage(c *evedata.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {

	p := NewPage(s, r, "EVE ECM Jam")
	templates.Templates = template.Must(template.ParseFiles("templates/ecmjam.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
