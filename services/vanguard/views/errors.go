package views

import (
	"html/template"
	"log"
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
)

func init() {
	vanguard.AddNotFoundHandler(notFoundPage)
}

func httpErrCode(w http.ResponseWriter, err error, code int) {
	if err != nil {
		log.Printf("http error %s", err)
	}
	http.Error(w, http.StatusText(code), code)
}

func httpErr(w http.ResponseWriter, err error) {
	log.Printf("http error %s", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	p := newPage(r, "Page Not Found")

	template := template.Must(template.ParseFiles("templates/error/notFound.html", LayoutPath))
	err := template.ExecuteTemplate(w, "base", p)
	if err != nil {
		httpErr(w, err)
		return
	}

	httpErrCode(w, nil, http.StatusNotFound)
}
