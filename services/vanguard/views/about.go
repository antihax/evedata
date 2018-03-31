package views

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/templates"
)

func init() {
	vanguard.AddRoute("account", "GET", "/about", aboutPage)
	vanguard.AddRoute("help", "GET", "/help", helpPage)
}

type page struct {
	Title    string
	Template string
}

var aboutPages = map[string]page{
	"":        page{Title: "About EVEData.org", Template: "evedata.html"},
	"evedata": page{Title: "About EVEData.org", Template: "evedata.html"},
	"privacy": page{Title: "EVEData.org Privacy Policy", Template: "privacy.html"},
	"terms":   page{Title: "EVEData.org Terms", Template: "terms.html"},
}
var helpPages = map[string]page{
	"shares":       page{Title: "Sharing Data", Template: "shares.html"},
	"integrations": page{Title: "Integrations", Template: "integrations.html"},
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	page, ok := aboutPages[r.FormValue("page")]
	if ok {
		renderStatic(w, r, "about", page)
	} else {
		httpErrCode(w, errors.New("not found"), http.StatusNotFound)
	}
	return
}

func helpPage(w http.ResponseWriter, r *http.Request) {
	page, ok := helpPages[r.FormValue("page")]
	if ok {
		renderStatic(w, r, "help", page)
	} else {
		httpErrCode(w, errors.New("not found"), http.StatusNotFound)
	}
	return
}

func renderStatic(w http.ResponseWriter, r *http.Request, area string, page page) {
	p := newPage(r, page.Title)
	setCache(w, 60*60*24)
	templates.Templates = template.Must(
		template.ParseFiles("templates/"+area+"/"+page.Template, "templates/"+area+"/base.html", templates.LayoutPath),
	)
	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		return
	}
	return
}
