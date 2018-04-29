package views

import (
	"errors"
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
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
	"":        {Title: "About EVEData.org", Template: "evedata.html"},
	"evedata": {Title: "About EVEData.org", Template: "evedata.html"},
	"privacy": {Title: "EVEData.org Privacy Policy", Template: "privacy.html"},
	"terms":   {Title: "EVEData.org Terms", Template: "terms.html"},
}
var helpPages = map[string]page{
	"shares":       {Title: "Sharing Data", Template: "shares.html"},
	"integrations": {Title: "Integrations", Template: "integrations.html"},
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
