package views

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("ContactSync", "GET", "/contactSync", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w,
			"contactSync.html",
			time.Hour*24*31,
			newPage(r, "Contact Copiers"))
	})
	vanguard.AddAuthRoute("ContactSync", "PUT", "/U/contactSync", apiAddContactSync)
	vanguard.AddAuthRoute("ContactSync", "GET", "/U/contactSync", apiGetContactSyncs)
	vanguard.AddAuthRoute("ContactSync", "DELETE", "/U/contactSync", apiDeleteContactSync)
}

func apiAddContactSync(w http.ResponseWriter, r *http.Request) {
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
	renderJSON(w, v, 0)
}

func apiDeleteContactSync(w http.ResponseWriter, r *http.Request) {
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
