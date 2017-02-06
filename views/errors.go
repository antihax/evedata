package views

import (
	"html/template"
	"net/http"

	"github.com/antihax/evedata/appContext"
	evedata "github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddNotFoundHandler(notFoundPage)
}

func notFoundPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	p := newPage(r, "Page Not Found")

	templates.Templates = template.Must(template.ParseFiles("templates/error/notFound.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNotFound, nil
}
