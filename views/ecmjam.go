package views

import (
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"
)

func init() {
	evedata.AddRoute(evedata.Route{"ecmjam", "GET", "/ecmjam", ecmjamPage})
}

// FindAgents generate a list of agents based on user input
func ecmjamPage(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	p := Page{
		Title: "ECM Jam",
	}

	templates.Templates = template.Must(template.ParseFiles("templates/ecmjam.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
