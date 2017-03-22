package views

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("iskPerLP", "GET", "/iskPerLP", iskPerLPPage)
	evedata.AddRoute("iskPerLPCorpss", "GET", "/J/iskPerLPCorps", iskPerLPCorps)
	evedata.AddRoute("iskPerLP", "GET", "/J/iskPerLP", iskPerLP)
}

func iskPerLPPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	p := newPage(r, "ISK Per Loyalty Point")

	templates.Templates = template.Must(template.ParseFiles("templates/iskPerLP.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)
	if err != nil {
		httpErr(w, err)
		return
	}
}

func iskPerLPCorps(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	v, err := models.GetISKPerLPCorporations()
	if err != nil {
		httpErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}

func iskPerLP(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*30)
	q := r.FormValue("corp")
	v, err := models.GetISKPerLP(q)
	if err != nil {
		httpErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}
