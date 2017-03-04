package views

import (
	"html/template"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("account", "GET", "/", mainPage)
}

func mainPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var err error

	setCache(w, 60*60*24)

	p := newPage(r, "EVE Online Intel Data")
	templates.Templates = template.Must(template.ParseFiles("templates/mainPage.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
