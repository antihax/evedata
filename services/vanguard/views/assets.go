package views

import (
	"encoding/json"
	"errors"
	"strconv"

	"html/template"
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/evedata/services/vanguard/templates"
)

func init() {
	vanguard.AddRoute("assets", "GET", "/assets", assetsPage)
	vanguard.AddAuthRoute("assets", "GET", "/U/assets", assetsAPI)
	vanguard.AddAuthRoute("assets", "GET", "/U/assetLocations", assetLocationsAPI)
	vanguard.AddAuthRoute("assets", "GET", "/U/assetCharacters", assetCharactersAPI)
}

func assetsPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	p := newPage(r, "Asset Information")
	templates.Templates = template.Must(template.ParseFiles("templates/assets.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}

func assetCharactersAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	setCache(w, 5*60)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for assets"), http.StatusUnauthorized)
		return
	}

	v, err := models.GetAssetCharacters(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}

func assetLocationsAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	setCache(w, 5*60)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for asset locations"), http.StatusUnauthorized)
		return
	}

	// Get arguments
	filterCharacterID := 0
	filter := r.FormValue("filterCharacterID")
	if filter != "" {
		filterCharacterID, err = strconv.Atoi(filter)
		if err != nil {
			httpErrCode(w, err, http.StatusNotFound)
			return
		}
	}

	v, err := models.GetAssetLocations(characterID, (int32)(filterCharacterID))
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}

func assetsAPI(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		locationID        int64
		filterCharacterID int32
	)

	setCache(w, 5*60)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for asset API"), http.StatusUnauthorized)
		return
	}

	// Get arguments
	filter := r.FormValue("filterCharacterID")
	if filter != "" {
		filterCharacterID64, err := strconv.ParseInt(filter, 10, 64)
		if err != nil {
			httpErrCode(w, err, http.StatusNotFound)
			return
		}
		filterCharacterID = int32(filterCharacterID64)
	}

	location := r.FormValue("locationID")
	if location != "" {
		locationID, err = strconv.ParseInt(location, 10, 64)
		if err != nil {
			httpErrCode(w, err, http.StatusNotFound)
			return
		}
	}

	v, err := models.GetAssets(characterID, filterCharacterID, locationID)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)
}
