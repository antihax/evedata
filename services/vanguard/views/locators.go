package views

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/evedata/services/vanguard/templates"
)

func init() {
	vanguard.AddRoute("locatorShares", "GET", "/locatorShares", locatorSharesPage)
	vanguard.AddAuthRoute("locatorShares", "GET", "/U/locatorShares", apiGetLocatorShares)
	vanguard.AddAuthRoute("locatorShares", "DELETE", "/U/locatorShares", apiDeleteLocatorShare)
	vanguard.AddAuthRoute("locatorShares", "POST", "/U/locatorShares", apiAddLocatorShare)

	vanguard.AddRoute("locators", "GET", "/locatorResponses", locatorResponsesPage)
	vanguard.AddAuthRoute("locators", "GET", "/U/locatorResponses", apiGetLocatorResponses)
}

func locatorSharesPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	p := newPage(r, "Account Information")
	templates.Templates = template.Must(template.ParseFiles("templates/locatorShare.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}

func apiDeleteLocatorShare(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	entity, err := strconv.ParseInt(r.FormValue("entityID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.DeleteLocatorShare(characterID, int32(entity)); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiAddLocatorShare(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	entity, err := strconv.ParseInt(r.FormValue("entityID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.AddLocatorShare(characterID, int32(entity)); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiGetLocatorShares(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	v, err := models.GetLocatorShares(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func locatorResponsesPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	p := newPage(r, "Account Information")
	templates.Templates = template.Must(template.ParseFiles("templates/locatorResponses.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}

func apiGetLocatorResponses(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())
	c := vanguard.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	info, err := getAccountInformation(c, s)
	if err != nil {
		httpErr(w, err)
		return
	}

	v, err := models.GetLocatorResponses(characterID, info.Cursor.CursorCharacterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}
