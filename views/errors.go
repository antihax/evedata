package views

import (
	"evedata/appContext"
	evedata "evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddNotFoundHandler(notFoundPage)
}

func notFoundPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "Page Not Found")

	templates.Templates = template.Must(template.ParseFiles("templates/error/notFound.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNotFound, nil
}
