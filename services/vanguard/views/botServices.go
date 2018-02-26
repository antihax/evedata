package views

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/evedata/services/vanguard/templates"
)

func init() {

	vanguard.AddRoute("botServices", "GET", "/botServices", botServicesPage)
	vanguard.AddAuthRoute("botServices", "GET", "/U/botServices", apiGetBotServices)
	vanguard.AddAuthRoute("botServices", "DELETE", "/U/botServices", apiDeleteBotService)
	vanguard.AddAuthRoute("botServices", "POST", "/U/botServicesDiscord", apiAddDiscordBotService)

	vanguard.AddRoute("botServices", "GET", "/botDetails", botDetailsPage)
	vanguard.AddAuthRoute("botServices", "GET", "/U/botDetails", apiGetBotDetails)

	vanguard.AddAuthRoute("botServices", "GET", "/U/entitiesWithRoles", apiGetEntitiesWithRoles)
}

func botServicesPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	p := newPage(r, "Integrations")
	templates.Templates = template.Must(template.ParseFiles("templates/botServices.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}

func apiDeleteBotService(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	botServerID, err := strconv.ParseInt(r.FormValue("botServerID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.DeleteService(characterID, int32(botServerID)); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiAddDiscordBotService(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())
	g := vanguard.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	// decode int to validate
	_, err := strconv.Atoi(r.FormValue("serverID"))
	if err != nil {
		httpErrCode(w, err, http.StatusBadRequest)
		return
	}

	entityID, err := strconv.ParseInt(r.FormValue("entityID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusBadRequest)
		return
	}

	if err = g.Conservator.Call("Conservator.VerifyDiscord", r.FormValue("serverID"), &ok); err != nil {
		httpErr(w, err)
		return
	}

	if !ok {
		httpErr(w, errors.New("serverID is invalid or the bot has no access."))
	}

	if err = models.AddDiscordService(characterID, int32(entityID), r.FormValue("serverID")); err != nil {
		httpErr(w, err)
		return
	}

	return
}

func apiGetBotServices(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	v, err := models.GetBotServices(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func botDetailsPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	p := newPage(r, "Integration Services")
	templates.Templates = template.Must(template.ParseFiles("templates/botDetails.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}
func apiGetBotDetails(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	serverID, err := strconv.Atoi(r.FormValue("serviceID"))
	if err != nil {
		httpErr(w, err)
		return
	}

	v, err := models.GetBotServiceDetails(characterID, int32(serverID))
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiGetEntitiesWithRoles(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	v, err := models.GetEntitiesWithRole(characterID, r.FormValue("role"))
	if err != nil {
		httpErr(w, err)
		return
	}

	json.NewEncoder(w).Encode(v)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}
