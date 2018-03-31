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

func aboutPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60*24)

	page, ok := aboutPages[r.FormValue("page")]
	if ok {
		p := newPage(r, page.Title)
		templates.Templates = template.Must(
			template.ParseFiles("templates/about/"+page.Template, "templates/about.html", templates.LayoutPath),
		)
		if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
			httpErrCode(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		httpErrCode(w, errors.New("not found"), http.StatusNotFound)
	}
	return
}
