package views

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var validEntity map[string]bool

func init() {
	vanguard.AddRoute("GET", "/alliance", alliancePage)
	vanguard.AddRoute("GET", "/character", characterPage)
	vanguard.AddRoute("GET", "/corporation", corporationPage)
	vanguard.AddRoute("GET", "/J/warsForEntity", warsForEntityAPI)
	vanguard.AddRoute("GET", "/J/shipsForEntity", shipsForEntityAPI)
	vanguard.AddRoute("GET", "/J/alliesForEntity", alliesForEntityAPI)
	vanguard.AddRoute("GET", "/J/heatmapForEntity", heatmapForEntityAPI)
	vanguard.AddRoute("GET", "/J/activityForEntity", activityForEntityAPI)
	vanguard.AddRoute("GET", "/J/killmailsForEntity", killmailsForEntityAPI)
	vanguard.AddRoute("GET", "/J/corporationsForAlliance", corporationsForAllianceAPI)

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

func killmailsForEntityAPI(w http.ResponseWriter, r *http.Request) {

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

	v, err := models.GetKillmailsForEntity(id, entityType)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Hour*1)
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

func entityBlurb(name, entType string, eff float64, kills, losses, capKills int64, plural bool) string {

	p := message.NewPrinter(language.English)
	if eff == 0 && kills > 1 {
		eff = 1
	}
	desc := name + " is a "
	if eff > 0.99 && kills > 100 {
		desc += "godlike"
	} else if eff > 0.9 && kills > 75 {
		desc += "deadly"
	} else if eff > 0.75 && kills > 50 {
		desc += "dangerous"
	} else if eff > 0.6 && kills > 25 {
		desc += "decent"
	} else if eff > 0.4 && kills > 25 {
		desc += "mediocre"
	} else if eff > 0.3 && kills > 25 {
		desc += "bad"
	} else if eff > 0.1 && kills > 25 {
		desc += "really bad"
	} else if eff <= 0.1 && kills > 25 {
		desc += "awful"
	} else if kills > 1 {
		desc += "useless"
	} else {
		desc += "carebear"
	}

	if plural {
		desc += " " + entType + " who are "
	} else {
		desc += " " + entType + " who is "
	}

	if kills+losses > 50000 {
		desc += "obnoxiously active"
	} else if kills+losses > 2500 {
		desc += "very active"
	} else if kills+losses > 2500 {
		desc += "very active"
	} else if kills+losses > 250 {
		desc += "quite active"
	} else if kills+losses > 50 {
		desc += "a little active"
	} else if kills+losses > 25 {
		desc += "not very active"
	} else {
		desc += "quitting eve"
	}

	capProb := float32(0)
	if capKills > 0 {
		capProb = float32(capKills) / float32(kills)
	}

	if capProb > 0.5 && capKills > 20 {
		desc += " and drop capital ships far too much."
	} else if capProb > 0.4 && capKills > 20 {
		desc += " and drop capital ships a lot."
	} else if capProb > 0.25 && capKills > 20 {
		desc += " and drop capital ships often."
	} else if capProb > 0.10 && capKills > 5 {
		desc += " and occasionally drop capital ships."
	} else if capProb > 0.01 && capKills > 5 {
		desc += " and rarely drop capital ships."
	} else {
		desc += "."
	}

	if kills+losses+capKills > 0 {
		desc += " In the last 90 days they have " + p.Sprintf("%d", kills) + " kills"
		if capKills > 0 {
			desc += p.Sprintf(", %d with capital ships involved,", capKills)
		}
		desc += " and " + p.Sprintf("%d", losses) + " losses."

	} else {
		desc += " They have no kill activity in the last 90 days."
	}

	return desc
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
	if ref.Kills > 0 && ref.Losses > 0 {
		ref.Efficiency = 1 - (float64(ref.Losses) / float64(ref.Kills))
	}

	description := entityBlurb(ref.AllianceName, "alliance", ref.Efficiency, ref.Kills, ref.Losses, ref.CapKills, true)
	p["OG"] = OpenGraph{
		Image:       entityImage(ref.AllianceID, "alliance", 128),
		Title:       "Alliance: " + ref.AllianceName + " - EveData.org",
		Description: description,
	}
	p["Description"] = description

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
	if ref.Kills > 0 && ref.Losses > 0 {
		ref.Efficiency = 1 - (float64(ref.Losses) / float64(ref.Kills))
	}

	description := entityBlurb(ref.CorporationName, "corporation", ref.Efficiency, ref.Kills, ref.Losses, ref.CapKills, true)
	if ref.AllianceID > 0 {
		description += " They are part of " + ref.AllianceName.String + "."
	}

	p["OG"] = OpenGraph{
		Image:       entityImage(ref.CorporationID, "corporation", 128),
		Title:       "Corporation: " + ref.CorporationName + " - EveData.org",
		Description: description,
	}
	p["Description"] = description
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

	if ref.Kills > 0 && ref.Losses > 0 {
		ref.Efficiency = 1 - (float64(ref.Losses) / float64(ref.Kills))
	}
	description := entityBlurb(ref.CharacterName, "character", ref.Efficiency, ref.Kills, ref.Losses, ref.CapKills, false)

	if ref.AllianceID > 0 {
		description += " They are part of " + ref.CorporationName + " with " + ref.AllianceName.String + "."
	} else {
		description += " They are part of " + ref.CorporationName + "."
	}

	p["OG"] = OpenGraph{
		Image:       entityImage(int64(ref.CharacterID), "character", 128),
		Title:       "Character: " + ref.CharacterName + " - EveData.org",
		Description: description,
	}
	p["Description"] = description
	p["Title"] = ref.CharacterName

	renderTemplate(w, "entities.html", time.Hour*24*7, p)
}
