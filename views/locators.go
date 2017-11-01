package views

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/evedata"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("locatorShares", "GET", "/locatorShares", locatorSharesPage)
	evedata.AddAuthRoute("locatorShares", "GET", "/U/locatorShares", apiGetLocatorShares)
	evedata.AddAuthRoute("locatorShares", "DELETE", "/U/locatorShares", apiDeleteLocatorShare)
	evedata.AddAuthRoute("locatorShares", "POST", "/U/locatorShares", apiAddLocatorShare)

	evedata.AddRoute("locators", "GET", "/locatorResponses", locatorResponsesPage)
	evedata.AddAuthRoute("locators", "GET", "/U/locatorResponses", apiGetLocatorResponses)
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
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	entity, err := strconv.ParseInt(r.FormValue("entityID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.DeleteLocatorShare(characterID, entity); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiAddLocatorShare(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	entity, err := strconv.ParseInt(r.FormValue("entityID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.AddLocatorShare(characterID, entity); err != nil {
		fmt.Println(err)
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiGetLocatorShares(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
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
	s := evedata.SessionFromContext(r.Context())
	c := evedata.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
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
