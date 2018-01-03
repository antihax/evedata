package views

import (
	"html/template"
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/templates"
)

func init() {
	vanguard.AddRoute("account", "GET", "/", mainPage)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	var err error

	setCache(w, 60*60*24)

	p := newPage(r, "EVE Online Intel Data")
	templates.Templates = template.Must(template.ParseFiles("templates/mainPage.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)
	if err != nil {
		httpErr(w, err)
		return
	}
}
