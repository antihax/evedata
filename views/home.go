package views

import (
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"
)

func init() {
	evedata.AddRoute(evedata.Route{"agents", "GET", "/", mainPage})
}

// FindAgents generate a list of agents based on user input
func mainPage(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	p := Page{
		Title: "home",
	}

	templates.Templates = template.Must(template.ParseFiles("templates/home.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
