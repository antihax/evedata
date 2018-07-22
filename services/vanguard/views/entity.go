package views

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

var validEntity map[string]bool

func init() {
	vanguard.AddRoute("GET", "/alliance", alliancePage)

	vanguard.AddRoute("GET", "/corporation", corporationPage)
	vanguard.AddRoute("GET", "/character", characterPage)
	vanguard.AddRoute("GET", "/J/warsForEntity", warsForEntityAPI)
	vanguard.AddRoute("GET", "/J/activityForEntity", activityForEntityAPI)
	vanguard.AddRoute("GET", "/J/heatmapForEntity", heatmapForEntityAPI)
	vanguard.AddRoute("GET", "/J/assetsForEntity", assetsForEntityAPI)
	vanguard.AddRoute("GET", "/J/alliesForEntity", alliesForEntityAPI)
	vanguard.AddRoute("GET", "/J/shipsForEntity", shipsForEntityAPI)
	vanguard.AddRoute("GET", "/J/corporationHistory", corporationHistoryAPI)
	vanguard.AddRoute("GET", "/J/corporationsForAlliance", corporationsForAllianceAPI)
	vanguard.AddRoute("GET", "/J/knownAssociatesForEntity", knownAssociatesForEntityAPI)
	vanguard.AddRoute("GET", "/J/allianceHistoryForEntity", allianceHistoryForEntityAPI)
	vanguard.AddRoute("GET", "/J/corporationHistoryForEntity", corporationHistoryForEntityAPI)
	vanguard.AddRoute("GET", "/J/allianceJoinHistoryForEntity", allianceJoinHistoryForEntityAPI)

	validEntity = map[string]bool{"alliance": true, "corporation": true, "character": true}
}

func shipsForEntityAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	entityType := r.FormValue("entityType")
	if !validEntity[entityType] {
		httpErr(w, errors.New("entityType must be corporation, character, or alliance"))
		return
	}

	v, err := models.GetKnownShipTypes(id, entityType)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func assetsForEntityAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	entityType := r.FormValue("entityType")
	if !validEntity[entityType] {
		httpErr(w, errors.New("entityType must be corporation, character, or alliance"))
		return
	}
	var v []models.AssetsInSpace
	if entityType == "alliance" {
		v, err = models.GetAllianceAssetsInSpace(id)
	} else if entityType == "corporation" {
		v, err = models.GetCorporationAssetsInSpace(id)
	} else {
		httpErr(w, errors.New("entityType must be corporation or alliance"))
		return
	}

	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func activityForEntityAPI(w http.ResponseWriter, r *http.Request) {

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	entityType := r.FormValue("entityType")
	if !validEntity[entityType] {
		httpErr(w, errors.New("entityType must be corporation, character, or alliance"))
		return
	}

	v, err := models.GetConstellationActivity(id, entityType)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func heatmapForEntityAPI(w http.ResponseWriter, r *http.Request) {

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	entityType := r.FormValue("entityType")
	if !validEntity[entityType] {
		httpErr(w, errors.New("entityType must be corporation, character, or alliance"))
		return
	}

	v, err := models.GetKillmailHeatMap(id, entityType)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func alliesForEntityAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	v, err := models.GetKnownAlliesByID(id)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func corporationHistoryAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	v, err := models.GetCorporationHistory(int32(id))
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func knownAssociatesForEntityAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	entityType := r.FormValue("entityType")
	if !validEntity[entityType] {
		httpErr(w, errors.New("entityType must be corporation, character, or alliance"))
		return
	}

	var v []models.KnownAlts
	if entityType == "alliance" {
		v, err = models.GetAllianceKnownAssociates(id)
	} else if entityType == "corporation" {
		v, err = models.GetCorporationKnownAssociates(id)
	} else if entityType == "character" {
		v, err = models.GetCharacterKnownAssociates(id)
	} else {
		httpErr(w, errors.New("entityType must be corporation or alliance"))
		return
	}

	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func allianceJoinHistoryForEntityAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	v, err := models.GetAllianceJoinHistory(id)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func allianceHistoryForEntityAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	v, err := models.GetAllianceHistory(id)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func corporationHistoryForEntityAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	v, err := models.GetCorporationJoinHistory(id)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func warsForEntityAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	v, err := models.GetWarsForEntityByID(id)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*4)
}

func corporationsForAllianceAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	v, err := models.GetAllianceMembers(id)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*12)
}

func alliancePage(w http.ResponseWriter, r *http.Request) {
	p := newPage(r, "Unknown Alliance")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	p["entityID"] = idStr
	p["entityType"] = "alliance"

	ref, err := models.GetAlliance(id)
	if err != nil {
		httpErr(w, err)
		return
	}
	p["Alliance"] = ref
	p["Title"] = ref.AllianceName

	renderTemplate(w, "entities.html", time.Hour*24*7, p)
}

func corporationPage(w http.ResponseWriter, r *http.Request) {
	p := newPage(r, "Unknown Corporation")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}
	p["entityID"] = idStr
	p["entityType"] = "corporation"

	ref, err := models.GetCorporation(id)
	if err != nil {
		httpErr(w, err)
		return
	}

	p["Corporation"] = ref
	p["Title"] = ref.CorporationName

	renderTemplate(w, "entities.html", time.Hour*24*7, p)
}

func characterPage(w http.ResponseWriter, r *http.Request) {
	p := newPage(r, "Unknown Character")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	p["entityID"] = idStr
	p["entityType"] = "character"

	ref, err := models.GetCharacter(int32(id))
	if err != nil {
		httpErr(w, err)
		return
	}

	p["Character"] = ref
	p["Title"] = ref.CharacterName

	renderTemplate(w, "entities.html", time.Hour*24*7, p)
}
