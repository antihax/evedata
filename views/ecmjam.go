package views

import (
	"html/template"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("ecmjam", "GET", "/ecmjam", ecmjamPage)
}

func ecmjamPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	p := newPage(r, "EVE ECM Jam")
	templates.Templates = template.Must(template.ParseFiles("templates/ecmjam.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
