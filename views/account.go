package views

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/eveapi"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("account", "GET", "/account", accountPage)

	evedata.AddAuthRoute("account", "GET", "/X/accountInfo", accountInfo)
	evedata.AddAuthRoute("account", "POST", "/X/cursorChar", cursorChar)

	evedata.AddAuthRoute("crestTokens", "GET", "/U/crestTokens", apiGetCRESTTokens)
	evedata.AddAuthRoute("crestTokens", "DELETE", "/U/crestTokens", apiDeleteCRESTToken)
}

func accountPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 0)
	p := newPage(r, "Account Information")
	templates.Templates = template.Must(template.ParseFiles("templates/account.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func accountInfo(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusForbidden, nil
	}

	char, ok := s.Values["character"].(eveapi.VerifyResponse)
	if !ok {
		return http.StatusForbidden, nil
	}

	accountInfo, ok := s.Values["accountInfo"].([]byte)
	if !ok {
		if err := updateAccountInfo(s, characterID, char.CharacterName); err != nil {
			return http.StatusInternalServerError, err
		}

		if err := s.Save(r, w); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	w.Write(accountInfo)

	return http.StatusOK, nil
}

func cursorChar(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusForbidden, nil
	}

	// Parse the cursorCharacterID
	cursorCharacterID, err := strconv.ParseInt(r.FormValue("cursorCharacterID"), 10, 64)
	if err != nil {
		return http.StatusForbidden, nil
	}

	// Set our new cursor
	err = models.SetCursorCharacter(characterID, cursorCharacterID)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	char, ok := s.Values["character"].(eveapi.VerifyResponse)
	if !ok {
		return http.StatusForbidden, nil
	}

	// Update the account information in redis
	if err = updateAccountInfo(s, characterID, char.CharacterName); err != nil {
		return http.StatusInternalServerError, err
	}

	if err = s.Save(r, w); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func apiGetCRESTTokens(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusForbidden, nil
	}

	tokens, err := models.GetCRESTTokens(characterID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(tokens)

	if err = s.Save(r, w); err != nil {
		return http.StatusInternalServerError, err
	}

	return 200, nil
}

func apiDeleteCRESTToken(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusForbidden, nil
	}

	cid, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		return http.StatusNotFound, errors.New("Invalid tokenCharacterID")
	}

	if err := models.DeleteCRESTToken(characterID, cid); err != nil {
		return http.StatusConflict, err
	}

	char, ok := s.Values["character"].(eveapi.VerifyResponse)
	if !ok {
		return http.StatusForbidden, nil
	}

	if err = updateAccountInfo(s, characterID, char.CharacterName); err != nil {
		return http.StatusInternalServerError, err
	}

	if err = s.Save(r, w); err != nil {
		return http.StatusInternalServerError, err
	}
	return 200, nil
}
