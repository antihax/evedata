package views

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
)

var validEntity map[string]bool

func init() {
	evedata.AddRoute("entity", "GET", "/alliance", alliancePage)

	evedata.AddRoute("entity", "GET", "/corporation", corporationPage)
	evedata.AddRoute("entity", "GET", "/character", characterPage)
	evedata.AddRoute("entity", "GET", "/J/warsForEntity", warsForEntityAPI)
	evedata.AddRoute("entity", "GET", "/J/activityForEntity", activityForEntityAPI)
	evedata.AddRoute("entity", "GET", "/J/assetsForEntity", assetsForEntityAPI)
	evedata.AddRoute("entity", "GET", "/J/alliesForEntity", alliesForEntityAPI)
	evedata.AddRoute("entity", "GET", "/J/shipsForEntity", shipsForEntityAPI)
	evedata.AddRoute("entity", "GET", "/J/corporationsForAlliance", corporationsForAllianceAPI)

	validEntity = map[string]bool{"alliance": true, "corporation": true, "character": true}
}

func shipsForEntityAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60*4)
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid id. Please provide a valid ?id=")
	}

	entityType := r.FormValue("entityType")
	if !validEntity[entityType] {
		return http.StatusInternalServerError, errors.New("entityType must be corporation, character, or alliance")
	}

	v, err := models.GetKnownShipTypes(id, entityType)
	if err != nil {
		return http.StatusNotFound, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
	return http.StatusOK, nil
}

func assetsForEntityAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60*4)
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid id. Please provide a valid ?id=")
	}

	entityType := r.FormValue("entityType")
	if !validEntity[entityType] {
		return http.StatusInternalServerError, errors.New("entityType must be corporation, character, or alliance")
	}
	var v []models.AssetsInSpace
	if entityType == "alliance" {
		v, err = models.GetAllianceAssetsInSpace(id)
	} else if entityType == "corporation" {
		v, err = models.GetCorporationAssetsInSpace(id)
	} else {
		return http.StatusInternalServerError, errors.New("entityType must be corporation or alliance")
	}

	if err != nil {
		return http.StatusNotFound, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
	return http.StatusOK, nil
}

func activityForEntityAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60*4)
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid id. Please provide a valid ?id=")
	}

	entityType := r.FormValue("entityType")
	if !validEntity[entityType] {
		return http.StatusInternalServerError, errors.New("entityType must be corporation, character, or alliance")
	}

	v, err := models.GetConstellationActivity(id, entityType)
	if err != nil {
		return http.StatusNotFound, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
	return http.StatusOK, nil
}

func alliesForEntityAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60*4)
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid id. Please provide a valid ?id=")
	}
	v, err := models.GetKnownAlliesByID(id)
	if err != nil {
		return http.StatusNotFound, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
	return http.StatusOK, nil
}

func warsForEntityAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid id. Please provide a valid ?id=")
	}
	v, err := models.GetWarsForEntityByID(id)
	if err != nil {
		return http.StatusNotFound, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
	return http.StatusOK, nil
}

func corporationsForAllianceAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid id. Please provide a valid ?id=")
	}
	v, err := models.GetAllianceMembers(id)
	if err != nil {
		return http.StatusNotFound, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
	return http.StatusOK, nil
}

func alliancePage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	p := newPage(r, "Unknown Alliance")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid Alliance ID. Please provide an ?id=")
	}

	p["entityID"] = idStr
	p["entityType"] = "alliance"

	ref, err := models.GetAlliance(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	p["Alliance"] = ref
	p["Title"] = ref.AllianceName

	templates.Templates = template.Must(template.ParseFiles("templates/entities.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func corporationPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	p := newPage(r, "Unknown Corporation")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid corporation ID. Please provide an ?id=")
	}
	p["entityID"] = idStr
	p["entityType"] = "corporation"

	ref, err := models.GetCorporation(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	p["Corporation"] = ref
	p["Title"] = ref.CorporationName

	templates.Templates = template.Must(template.ParseFiles("templates/entities.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func characterPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	p := newPage(r, "Unknown Character")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid character ID. Please provide an ?id=")
	}

	p["entityID"] = idStr
	p["entityType"] = "character"

	ref, err := models.GetCharacter(id)
	if err != nil {

		return http.StatusInternalServerError, err
	}
	p["Character"] = ref
	p["Title"] = ref.CharacterName

	templates.Templates = template.Must(template.ParseFiles("templates/entities.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
