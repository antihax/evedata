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
	vanguard.AddRoute("ContactSync", "GET", "/contactSync", contactSyncPage)
	vanguard.AddAuthRoute("ContactSync", "PUT", "/U/contactSync", apiAddContactSync)
	vanguard.AddAuthRoute("ContactSync", "GET", "/U/contactSync", apiGetContactSyncs)
	vanguard.AddAuthRoute("ContactSync", "DELETE", "/U/contactSync", apiDeleteContactSync)
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
	s := vanguard.SessionFromContext(r.Context())

	type localContactSync struct {
		Source      int `json:",string"`
		Destination int `json:",string"`
	}
	var cc localContactSync

	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	if r.Body == nil {
		httpErrCode(w, nil, http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&cc)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.AddContactSync(characterID, cc.Source, cc.Destination); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiGetContactSyncs(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	v, err := models.GetContactSyncs(characterID)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(v)
}

func apiDeleteContactSync(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	destination, err := strconv.Atoi(r.FormValue("destination"))
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.DeleteContactSync(characterID, destination); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}
