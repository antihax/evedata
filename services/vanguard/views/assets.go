package views

import (
	"errors"
	"log"
	"strconv"
	"time"

	"net/http"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"

	"github.com/antihax/goesi"
)

func init() {
	vanguard.AddRoute("GET", "/assets", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "assets.html", time.Hour*24*31, newPage(r, "Asset Information"))
	})
	vanguard.AddRoute("GET", "/marketableAssets", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "marketableAssets.html", time.Hour*24*31, newPage(r, "Marketable Asset Valuation"))
	})
	vanguard.AddAuthRoute("GET", "/U/assets", assetsAPI)
	vanguard.AddAuthRoute("GET", "/U/marketableAssets", marketableAssetsAPI)
	vanguard.AddAuthRoute("GET", "/U/assetLocations", assetLocationsAPI)
	vanguard.AddAuthRoute("GET", "/U/assetCharacters", assetCharactersAPI)
}

func assetCharactersAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	ch, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		log.Println(err)
		httpErrCode(w, errors.New("could not find character ID for assets"), http.StatusUnauthorized)
		return
	}

	marketable := false
	if r.FormValue("marketable") == "1" {
		marketable = true
	}

	v, err := models.GetAssetCharacters(ch.CharacterID, ch.CharacterOwnerHash, marketable)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	if len(v) == 0 {
		httpErrCode(w, errors.New("No asset characters"), http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour)
}

func assetLocationsAPI(w http.ResponseWriter, r *http.Request) {
	var err error

	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	ch, ok := s.Values["character"].(goesi.VerifyResponse)
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
			log.Println(err)
			httpErrCode(w, err, http.StatusNotFound)
			return
		}
	}
	marketable := false
	if r.FormValue("marketable") == "1" {
		marketable = true
	}
	v, err := models.GetAssetLocations(ch.CharacterID, ch.CharacterOwnerHash, (int32)(filterCharacterID), marketable)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	if len(v) == 0 {
		httpErrCode(w, errors.New("No asset locations"), http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour)
}

func assetsAPI(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		locationID        int64
		filterCharacterID int32
	)

	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	ch, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for asset API"), http.StatusUnauthorized)
		return
	}

	// Get arguments
	filter := r.FormValue("filterCharacterID")
	if filter != "" {
		filterCharacterID64, err := strconv.ParseInt(filter, 10, 64)
		if err != nil {
			log.Println(err)
			httpErrCode(w, err, http.StatusNotFound)
			return
		}
		filterCharacterID = int32(filterCharacterID64)
	}

	location := r.FormValue("locationID")
	if location != "" {
		locationID, err = strconv.ParseInt(location, 10, 64)
		if err != nil {
			log.Println(err)
			httpErrCode(w, err, http.StatusNotFound)
			return
		}
	}

	v, err := models.GetAssets(ch.CharacterID, ch.CharacterOwnerHash, filterCharacterID, locationID)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}

func marketableAssetsAPI(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	ch, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for marketable asset API"), http.StatusUnauthorized)
		return
	}

	// Get arguments
	tokenCharacterID, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		log.Println(err)
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	locationID, err := strconv.ParseInt(r.FormValue("locationID"), 10, 64)
	if err != nil {
		log.Println(err)
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	// Get Assets
	v, err := models.GetMarketableAssets(ch.CharacterID, ch.CharacterOwnerHash, int32(tokenCharacterID), locationID)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}
