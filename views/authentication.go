package views

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddAuthRoute("logout", "GET", "/X/logout", logout)

	evedata.AddAuthRoute("eveAuth", "GET", "/X/eveAuth", eveSSO)
	evedata.AddAuthRoute("eveSSOAnswer", "GET", "/X/eveSSOAnswer", eveSSOAnswer)

	evedata.AddAuthRoute("eveTokenAuth", "GET", "/X/eveTokenAuth", eveCRESTToken)
	evedata.AddAuthRoute("eveTokenAnswer", "GET", "/X/eveTokenAnswer", eveTokenAnswer)
}

func logout(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	s.Options.MaxAge = -1
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/", 302)
	return http.StatusMovedPermanently, nil
}

func eveSSO(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	s.Values["state"] = state

	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	url := c.SSOAuthenticator.AuthorizeURL(state, true, nil)
	http.Redirect(w, r, url, 302)
	return http.StatusMovedPermanently, nil
}

func eveSSOAnswer(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["state"] != state {
		return http.StatusInternalServerError, errors.New("State does not match. We likely could not read the sessin cookie. Please make sure cookies are enabled.")
	}

	tok, err := c.SSOAuthenticator.TokenExchange(code)
	if err != nil {
		return http.StatusInternalServerError, errors.New(fmt.Sprintf("%s code was %s and state was %s\n", err, code, state))
	}

	tokSrc, err := c.SSOAuthenticator.TokenSource(tok)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	v, err := c.EVE.Verify(tokSrc)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	s.Values["character"] = v
	s.Values["characterID"] = v.CharacterID
	s.Values["token"] = tok

	if err = updateAccountInfo(s, v.CharacterID, v.CharacterName); err != nil {
		return http.StatusInternalServerError, err
	}

	if err = s.Save(r, w); err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/account", 302)
	return http.StatusMovedPermanently, nil
}

type accountInformation struct {
	CharacterID   int64                  `json:"characterID"`
	CharacterName string                 `json:"characterName"`
	Characters    []models.CRESTToken    `json:"characters"`
	Cursor        models.CursorCharacter `json:"cursor"`
}

func updateAccountInfo(s *sessions.Session, characterID int64, characterName string) error {
	var err error
	a := accountInformation{}

	a.CharacterName = characterName
	a.CharacterID = characterID
	a.Characters, err = models.GetCRESTTokens(characterID)
	if err != nil {
		return err
	}

	a.Cursor, err = models.GetCursorCharacter(characterID)
	b, err := json.Marshal(a)
	s.Values["accountInfo"] = b

	return err
}

func eveCRESTToken(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)

	// Get the scopeGroups
	scopeGroupsTxt := r.FormValue("scopeGroups")
	if scopeGroupsTxt == "" {
		return http.StatusBadRequest, errors.New("scopeGroups is empty")
	}

	// split into []string
	scopeGroups := strings.Split(scopeGroupsTxt, ",")

	// Validate the scopeGroups are actually real
	validate := models.GetCharacterScopeGroups()
	for _, group := range scopeGroups {
		if validate[group] == "" {
			return http.StatusBadRequest, errors.New("scopeGroup is invalid")
		}
	}

	// Get the associated scopes to the groups
	scopes := models.GetCharacterScopesByGroups(scopeGroups)

	// Make a code to validate on the return
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Save the code to our session store to compare later
	s.Values["TOKENstate"] = state
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	url := c.TokenAuthenticator.AuthorizeURL(state, true, scopes)

	http.Redirect(w, r, url, 302)
	return http.StatusMovedPermanently, nil
}

func eveTokenAnswer(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["TOKENstate"] != state {
		return http.StatusInternalServerError, errors.New("Invalid State. It is possible that the session cookie is missing. Stop eating the cookies!")
	}

	tok, err := c.TokenAuthenticator.TokenExchange(code)
	if err != nil {
		return http.StatusInternalServerError, errors.New(fmt.Sprintf("%s code was %s and state was %s\n", err, code, state))
	}

	tokSrc, err := c.SSOAuthenticator.TokenSource(tok)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	v, err := c.EVE.Verify(tokSrc)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	characterID := s.Values["characterID"].(int64)
	err = models.AddCRESTToken(characterID, v.CharacterID, v.CharacterName, tok, v.Scopes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/account", 302)
	return http.StatusMovedPermanently, nil
}
