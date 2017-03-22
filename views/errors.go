package views

import (
	"html/template"
	"net/http"

	evedata "github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddNotFoundHandler(notFoundPage)
}

func httpErrCode(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func httpErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	p := newPage(r, "Page Not Found")

	templates.Templates = template.Must(template.ParseFiles("templates/error/notFound.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)
	if err != nil {
		httpErr(w, err)
		return
	}

	httpErrCode(w, http.StatusNotFound)
}
