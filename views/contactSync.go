package views

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/evedata"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("ContactSync", "GET", "/contactSync", contactSyncPage)
	evedata.AddAuthRoute("ContactSync", "PUT", "/U/contactSync", apiAddContactSync)
	evedata.AddAuthRoute("ContactSync", "GET", "/U/contactSync", apiGetContactSyncs)
	evedata.AddAuthRoute("ContactSync", "DELETE", "/U/contactSync", apiDeleteContactSync)
}

func contactSyncPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	p := newPage(r, "Contact Copiers")
	templates.Templates = template.Must(template.ParseFiles("templates/contactSync.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}

func apiAddContactSync(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	type localContactSync struct {
		Source      int `json:",string"`
		Destination int `json:",string"`
	}
	var cc localContactSync

	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	if r.Body == nil {
		httpErrCode(w, http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&cc)
	if err != nil {
		httpErrCode(w, http.StatusNotFound)
		return
	}

	if err := models.AddContactSync(characterID, cc.Source, cc.Destination); err != nil {
		httpErrCode(w, http.StatusConflict)
		return
	}
}

func apiGetContactSyncs(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	cc, err := models.GetContactSyncs(characterID)
	if err != nil {
		httpErrCode(w, http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(cc)
}

func apiDeleteContactSync(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	destination, err := strconv.Atoi(r.FormValue("destination"))
	if err != nil {
		httpErrCode(w, http.StatusNotFound)
		return
	}

	if err := models.DeleteContactSync(characterID, destination); err != nil {
		httpErrCode(w, http.StatusConflict)
		return
	}
}
