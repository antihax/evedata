package views

import (
	"encoding/json"
	"errors"
	"evedata/appContext"
	"evedata/models"
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("account", "GET", "/account", accountPage)

	evedata.AddRoute("crestTokens", "GET", "/U/crestTokens", apiGetCRESTTokens)
	evedata.AddRoute("crestTokens", "DELETE", "/U/crestTokens", apiDeleteCRESTToken)
}

func accountPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {

	p := newPage(s, r, "Account Information")
	templates.Templates = template.Must(template.ParseFiles("templates/account.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func apiGetCRESTTokens(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	if s.Values["characterID"] == nil {
		return http.StatusForbidden, nil
	}
	characterID := s.Values["characterID"].(int64)

	tokens, err := models.GetCRESTTokens(characterID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(tokens)

	return 200, nil
}

func apiDeleteCRESTToken(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {

	cid, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		return http.StatusNotFound, errors.New("Invalid tokenCharacterID")
	}
	characterID := s.Values["characterID"].(int64)
	if err := models.DeleteCRESTToken(characterID, cid); err != nil {
		return http.StatusConflict, err
	}

	return 200, nil
}
