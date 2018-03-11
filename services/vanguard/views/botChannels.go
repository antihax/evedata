package views

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddAuthRoute("botServices", "GET", "/U/botServiceChannels", apiGetBotServiceChannels)
	vanguard.AddAuthRoute("botServices", "POST", "/U/botServiceChannels", apiAddBotServiceChannel)
	vanguard.AddAuthRoute("botServices", "DELETE", "/U/botServiceChannels", apiDeleteBotServiceChannel)
}

func apiGetBotServiceChannels(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	g := vanguard.GlobalsFromContext(r.Context())

	// Verify the user has access to this service
	service, err := getBotService(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	channels := [][]string{}
	if err := g.RPCall("Conservator.GetChannels", service.BotServiceID, &channels); err != nil {
		httpErr(w, err)
		return
	}

	type channel struct {
		ChannelID   string `json:"channelID"`
		ChannelName string `json:"channelName"`
	}
	convChannels := []channel{}

	for _, ch := range channels {
		convChannels = append(convChannels, channel{ChannelID: ch[0], ChannelName: ch[1]})
	}
	json.NewEncoder(w).Encode(convChannels)
}

func apiAddBotServiceChannel(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)

	g := vanguard.GlobalsFromContext(r.Context())

	// Verify the user has access to this service
	service, err := getBotService(r)
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
			if err := g.RPCall("Conservator.GetChannels", service.BotServiceID, &channels); err != nil {
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

	if err = models.AddBotServiceChannel(service.BotServiceID, channelID, channelName); err != nil {
		httpErr(w, err)
		return
	}

	return
}

func apiDeleteBotServiceChannel(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	// Verify the user has access to this service
	service, err := getBotService(r)
	if err != nil {
		httpErr(w, err)
		return
	}

	channelID := r.FormValue("channelID")
	if channelID == "" {
		httpErrCode(w, nil, http.StatusTeapot)
	}

	if err := models.DeleteBotServiceChannel(service.BotServiceID, channelID); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}
