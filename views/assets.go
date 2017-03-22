package views

import (
	"encoding/json"
	"strconv"

	"html/template"
	"net/http"

	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("assets", "GET", "/assets", assetsPage)
	evedata.AddAuthRoute("assets", "GET", "/U/assets", assetsAPI)
	evedata.AddAuthRoute("assets", "GET", "/U/assetLocations", assetLocationsAPI)
	evedata.AddAuthRoute("assets", "GET", "/U/assetCharacters", assetCharactersAPI)
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
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	assetCharacters, err := models.GetAssetCharacters(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(assetCharacters)
}

func assetLocationsAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	setCache(w, 5*60)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	// Get arguments
	filterCharacterID := 0
	filter := r.FormValue("filterCharacterID")
	if filter != "" {
		filterCharacterID, err = strconv.Atoi(filter)
		if err != nil {
			httpErrCode(w, http.StatusNotFound)
			return
		}
	}

	assetLocations, err := models.GetAssetLocations(characterID, (int64)(filterCharacterID))
	if err != nil {
		httpErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(assetLocations)
}

func assetsAPI(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		locationID        int64
		filterCharacterID int64
	)

	setCache(w, 5*60)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	// Get arguments
	filter := r.FormValue("filterCharacterID")
	if filter != "" {
		filterCharacterID, err = strconv.ParseInt(filter, 10, 64)
		if err != nil {
			httpErrCode(w, http.StatusNotFound)
			return
		}
	}

	location := r.FormValue("locationID")
	if location != "" {
		locationID, err = strconv.ParseInt(location, 10, 64)
		if err != nil {
			httpErrCode(w, http.StatusNotFound)
			return
		}
	}

	assets, err := models.GetAssets(characterID, filterCharacterID, locationID)
	if err != nil {
		httpErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(assets)
}
