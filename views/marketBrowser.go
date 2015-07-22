package views

import (
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"
)

func init() {
	evedata.AddRoute(evedata.Route{"marketBrowser", "GET", "/", marketBrowser})
}

// marketBrowser generates.... stuff
func marketBrowser(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	p := Page{
		Title: "EVE Online Market Browser",
	}

	templates.Templates = template.Must(template.ParseFiles("templates/marketBrowser.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
