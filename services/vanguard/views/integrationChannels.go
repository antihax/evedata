package views

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/antihax/evedata/services/conservator"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddAuthRoute("GET", "/U/integrationChannels", apiGetIntegrationChannels)
	vanguard.AddAuthRoute("POST", "/U/integrationChannels", apiAddIntegrationChannel)
	vanguard.AddAuthRoute("DELETE", "/U/integrationChannels", apiDeleteIntegrationChannel)

	vanguard.AddAuthRoute("PUT", "/U/integrationChannelOptions", apiSetIntegrationChannelOptions)

	vanguard.AddAuthRoute("GET", "/U/integrationRoles", apiGetIntegrationRoles)
}

func apiGetIntegrationChannels(w http.ResponseWriter, r *http.Request) {
	g := vanguard.GlobalsFromContext(r.Context())

	// Verify the user has access to this service
	service, err := getIntegration(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	channels := [][]string{}
	if err := g.RPCall("Conservator.GetChannels", service.IntegrationID, &channels); err != nil {
		httpErr(w, err)
		return
	}

	type channel struct {
		ChannelID   string `json:"channelID"`
		ChannelName string `json:"channelName"`
	}
	v := []channel{}

	for _, ch := range channels {
		v = append(v, channel{ChannelID: ch[0], ChannelName: ch[1]})
	}
	renderJSON(w, v, time.Hour)
}

func apiGetIntegrationRoles(w http.ResponseWriter, r *http.Request) {
	g := vanguard.GlobalsFromContext(r.Context())

	// Verify the user has access to this service
	service, err := getIntegration(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	roles := [][]string{}
	if err := g.RPCall("Conservator.GetRoles", service.IntegrationID, &roles); err != nil {
		httpErr(w, err)
		return
	}

	type role struct {
		RoleID   string `json:"roleID"`
		RoleName string `json:"roleName"`
	}
	v := []role{}

	for _, ch := range roles {
		v = append(v, role{RoleID: ch[0], RoleName: ch[1]})
	}
	renderJSON(w, v, time.Hour)
}

func apiAddIntegrationChannel(w http.ResponseWriter, r *http.Request) {
	g := vanguard.GlobalsFromContext(r.Context())

	// Verify the user has access to this service
	service, err := getIntegration(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	ok := false
	channelName := ""
	channelID := r.FormValue("channelID")
	if service.Type == "discord" {
		// Verify the discord exists
		if err = g.RPCall("Conservator.VerifyDiscordChannel", []string{service.Address, channelID}, &ok); err != nil {
			httpErr(w, err)
			return
		}
		if ok {
			channels := [][]string{}
			if err := g.RPCall("Conservator.GetChannels", service.IntegrationID, &channels); err != nil {
				httpErr(w, err)
				return
			}
			for _, ch := range channels {
				if ch[0] == channelID {
					channelName = ch[1]
					break
				}
			}
		}
	}

	if !ok {
		httpErr(w, errors.New("serverID is invalid or the bot has no access."))
	}

	if err = models.AddIntegrationChannel(service.IntegrationID, channelID, channelName); err != nil {
		httpErr(w, err)
		return
	}
}

func apiDeleteIntegrationChannel(w http.ResponseWriter, r *http.Request) {
	// Verify the user has access to this service
	service, err := getIntegration(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	channelID := r.FormValue("channelID")
	if channelID == "" {
		httpErrCode(w, nil, http.StatusTeapot)
	}

	if err := models.DeleteIntegrationChannel(service.IntegrationID, channelID); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func apiSetIntegrationChannelOptions(w http.ResponseWriter, r *http.Request) {
	// Verify the user has access to this service
	service, err := getIntegration(r)
	if err != nil {
		httpErr(w, err)
		return
	}
	channelID := r.FormValue("channelID")
	if channelID == "" {
		httpErrCode(w, nil, http.StatusTeapot)
	}

	// unmarshal to verify this is accurate.
	chanOpts := conservator.ChannelOptions{}
	if err := json.Unmarshal([]byte(r.FormValue("options")), &chanOpts); err != nil {
		httpErr(w, err)
		return
	}

	//Unmarshal and format to our set string
	chanServices := conservator.ChannelTypes{}
	if err := json.Unmarshal([]byte(r.FormValue("services")), &chanServices); err != nil {
		httpErr(w, err)
		return
	}

	options, err := json.Marshal(chanOpts)
	if err != nil {
		httpErr(w, err)
		return
	}

	if err := models.UpdateChannel(service.IntegrationID, channelID, string(options), chanServices.GetServices()); err != nil {
		httpErr(w, err)
		return
	}
}
