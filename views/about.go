package views

import (
	"html/template"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("account", "GET", "/about", aboutPage)
}

func aboutPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var err error

	setCache(w, 60*60*24)

	page := r.FormValue("page")

	if page == "evedata" || page == "" {
		p := newPage(r, "About EVEData.org")
		templates.Templates = template.Must(template.ParseFiles("templates/about/evedata.html", "templates/about.html", templates.LayoutPath))
		err = templates.Templates.ExecuteTemplate(w, "base", p)
	} else if page == "privacy" {
		p := newPage(r, "EVEData.org Privacy Policy")
		templates.Templates = template.Must(template.ParseFiles("templates/about/privacy.html", "templates/about.html", templates.LayoutPath))
		err = templates.Templates.ExecuteTemplate(w, "base", p)
	} else if page == "terms" {
		p := newPage(r, "EVEData.org Terms")
		templates.Templates = template.Must(template.ParseFiles("templates/about/terms.html", "templates/about.html", templates.LayoutPath))
		err = templates.Templates.ExecuteTemplate(w, "base", p)
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
