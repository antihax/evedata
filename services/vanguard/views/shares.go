package views

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {

	vanguard.AddRoute("GET", "/shares",
		func(w http.ResponseWriter, r *http.Request) {
			p := newPage(r, "Share data with entities")
			p["ShareGroups"] = models.GetCharacterShareGroups()
			renderTemplate(w,
				"shares.html",
				time.Hour*24*31,
				p)
		})

	vanguard.AddAuthRoute("GET", "/U/shares", apiGetShares)
	vanguard.AddAuthRoute("DELETE", "/U/shares", apiDeleteShare)
	vanguard.AddAuthRoute("POST", "/U/shares", apiAddShare)
}

func apiDeleteShare(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	entityID, err := strconv.ParseInt(r.FormValue("entityID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	tokenCharacterID, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.DeleteShare(characterID, int32(tokenCharacterID), int32(entityID)); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiAddShare(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	entityID, err := strconv.ParseInt(r.FormValue("entityID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusBadRequest)
		return
	}

	tokenCharacterID, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusBadRequest)
		return
	}

	// verify these are real types
	types := strings.Split(r.FormValue("types"), ",")
	for _, t := range types {
		found := false
		for group := range models.GetCharacterShareGroups() {
			if t == group {
				found = true
				break
			}
		}
		if !found {
			httpErrCode(w, errors.New("Invalid Type"), http.StatusBadRequest)
			return
		}
	}

	if err := models.AddShare(characterID, int32(tokenCharacterID), int32(entityID), r.FormValue("types")); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiGetShares(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	v, err := models.GetShares(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, 0)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}
