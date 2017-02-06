package views

import (
	"encoding/json"
	"errors"
	"strconv"

	"html/template"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("assets", "GET", "/assets", assetsPage)
	evedata.AddAuthRoute("assets", "GET", "/U/assets", assetsAPI)
	evedata.AddAuthRoute("assets", "GET", "/U/assetLocations", assetLocationsAPI)
	evedata.AddAuthRoute("assets", "GET", "/U/assetCharacters", assetCharactersAPI)
}

func assetsPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	p := newPage(r, "Asset Information")
	templates.Templates = template.Must(template.ParseFiles("templates/assets.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func assetCharactersAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	var err error
	setCache(w, 5*60)

	if s.Values["characterID"] == nil || s.Values["characterID"] == 0 {
		return http.StatusForbidden, nil
	}

	// get our character ID from the session
	characterID := s.Values["characterID"].(int64)

	assetCharacters, err := models.GetAssetCharacters(characterID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(assetCharacters)

	return 200, nil
}

func assetLocationsAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	var err error
	setCache(w, 5*60)

	if s.Values["characterID"] == nil || s.Values["characterID"] == 0 {
		return http.StatusForbidden, nil
	}

	// Get arguments
	filterCharacterID := 0
	filter := r.FormValue("filterCharacterID")
	if filter != "" {
		filterCharacterID, err = strconv.Atoi(filter)
		if err != nil {
			return http.StatusNotFound, errors.New("Invalid filterCharacterID")
		}
	}

	// get our character ID from the session
	characterID := s.Values["characterID"].(int64)

	assetLocations, err := models.GetAssetLocations(characterID, (int64)(filterCharacterID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(assetLocations)

	return 200, nil
}

func assetsAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	var (
		err               error
		locationID        int64
		filterCharacterID int64
	)

	setCache(w, 5*60)

	if s.Values["characterID"] == nil || s.Values["characterID"] == 0 {
		return http.StatusForbidden, nil
	}
	characterID := s.Values["characterID"].(int64)

	// Get arguments

	filter := r.FormValue("filterCharacterID")
	if filter != "" {
		filterCharacterID, err = strconv.ParseInt(filter, 10, 64)
		if err != nil {
			return http.StatusNotFound, errors.New("Invalid filterCharacterID")
		}
	}

	location := r.FormValue("locationID")
	if location != "" {
		locationID, err = strconv.ParseInt(location, 10, 64)
		if err != nil {
			return http.StatusNotFound, errors.New("Invalid locationID")
		}
	}

	assets, err := models.GetAssets(characterID, filterCharacterID, locationID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(assets)

	return 200, nil
}
