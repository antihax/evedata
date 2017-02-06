package views

import (
	"github.com/antihax/evedata/appContext"
	evedata "github.com/antihax/evedata/server"

	"html/template"
	"net/http"

	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("top", "GET", "/top", topPage)
}

func topPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60*24*7)
	p := newPage(r, "EVEData.org backend statistics")
	templates.Templates = template.Must(template.ParseFiles("templates/top.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
