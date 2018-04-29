package views

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/conservator"
	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	// Integrations
	vanguard.AddRoute("integrations", "GET", "/integrations",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"integrations.html",
				time.Hour*24*31,
				newPage(r, "Integrations"))
		})

	vanguard.AddAuthRoute("integrations", "GET", "/U/integrations", apiGetIntegrations)
	vanguard.AddAuthRoute("integrations", "DELETE", "/U/integrations", apiDeleteIntegration)
	vanguard.AddAuthRoute("integrations", "POST", "/U/integrationsDiscord", apiAddDiscordIntegration)
	vanguard.AddAuthRoute("integrations", "POST", "/U/integrationShareToggleIgnore", apiIntegrationToggleIgnore)

	// Integration Details
	vanguard.AddRoute("integrations", "GET", "/integrationDetails", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w,
			"integrationDetails.html",
			time.Hour*24*31,
			newPage(r, "Integration Services"))
	})

	vanguard.AddAuthRoute("integrations", "GET", "/U/integrationDetails", apiGetIntegrationDetails)
	vanguard.AddAuthRoute("integrations", "PUT", "/U/integrationDetails", apiIntegrationOptions)

	vanguard.AddAuthRoute("integrations", "GET", "/U/entitiesWithRoles", apiGetEntitiesWithRoles)
}

func apiDeleteIntegration(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	integrationID, err := strconv.ParseInt(r.FormValue("integrationID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.DeleteService(characterID, int32(integrationID)); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiAddDiscordIntegration(w http.ResponseWriter, r *http.Request) {
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

	// Verify the discord exists
	if err = g.RPCall("Conservator.VerifyDiscord", r.FormValue("serverID"), &ok); err != nil {
		httpErr(w, err)
		return
	}

	if !ok {
		httpErr(w, errors.New("serverID is invalid or the bot has no access."))
		return
	}

	if err = models.AddDiscordService(characterID, int32(entityID), r.FormValue("serverID")); err != nil {
		httpErr(w, err)
		return
	}

	return
}

func apiGetIntegrations(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	v, err := models.GetIntegrations(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}
	renderJSON(w, v, 0)
}

func apiGetIntegrationDetails(w http.ResponseWriter, r *http.Request) {
	// Verify the user has access to this service
	v, err := getIntegration(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	json.Unmarshal([]byte(v.OptionsJSON), &v.Options)
	for i := range v.Channels {
		json.Unmarshal([]byte(v.Channels[i].OptionsJSON), &v.Channels[i].Options)
	}
	renderJSON(w, v, 0)
}

func apiIntegrationToggleIgnore(w http.ResponseWriter, r *http.Request) {
	g := vanguard.GlobalsFromContext(r.Context())

	// Check tokenCharacterID is valid
	tokenCharacterID, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		return
	}

	// Verify the user has access to this service
	service, err := getIntegration(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	_, err = g.Db.Exec("UPDATE evedata.sharing SET ignored = ! ignored WHERE entityID = ? AND tokenCharacterID = ?", service.EntityID, tokenCharacterID)
	if err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		return
	}
}

func apiGetEntitiesWithRoles(w http.ResponseWriter, r *http.Request) {
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

	renderJSON(w, v, 0)
}

func apiIntegrationOptions(w http.ResponseWriter, r *http.Request) {
	// Verify the user has access to this service
	service, err := getIntegration(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	// unmarshal to verify this is accurate.
	servOpts := conservator.ServiceOptions{}
	if err := json.Unmarshal([]byte(r.FormValue("options")), &servOpts); err != nil {
		httpErr(w, err)
		return
	}

	//Unmarshal and format to our set string
	servServices := conservator.ServiceTypes{}
	if err := json.Unmarshal([]byte(r.FormValue("services")), &servServices); err != nil {
		httpErr(w, err)
		return
	}

	options, err := json.Marshal(servOpts)
	if err != nil {
		httpErr(w, err)
		return
	}

	if err := models.UpdateService(service.IntegrationID, string(options), servServices.GetServices()); err != nil {
		httpErr(w, err)
		return
	}
}

func getIntegration(r *http.Request) (*models.IntegrationDetails, error) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		return nil, errors.New("Not authorized")
	}

	// Check integrationID is valid
	integrationID, err := strconv.Atoi(r.FormValue("integrationID"))
	if err != nil {
		return nil, err
	}

	// verify this character can access this service
	v, err := models.GetIntegrationDetails(characterID, int32(integrationID))
	if err != nil {
		return nil, err
	}

	if v.IntegrationID == 0 {
		return nil, errors.New("Not authorized")
	}
	return &v, nil
}
