package views

import (
	"encoding/json"
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
	vanguard.AddAuthRoute("botServices", "POST", "/U/botServices", apiAddBotService)
}

func botServicesPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	p := newPage(r, "Account Information")
	templates.Templates = template.Must(template.ParseFiles("templates/botServices.html", templates.LayoutPath))

	p["ShareGroups"] = models.GetCharacterShareGroups()

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

func apiAddBotService(w http.ResponseWriter, r *http.Request) {

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
