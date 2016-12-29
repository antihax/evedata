package views

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("logout", "GET", "/logout", logout)

	evedata.AddRoute("eveAuth", "GET", "/eveAuth", eveSSO)
	evedata.AddRoute("eveSSOAnswer", "GET", "/eveSSOAnswer", eveSSOAnswer)

	evedata.AddRoute("eveTokenAuth", "GET", "/eveTokenAuth", eveCRESTToken)
	evedata.AddRoute("eveTokenAnswer", "GET", "/eveTokenAnswer", eveTokenAnswer)
}

func logout(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	s.Options.MaxAge = -1
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/", 302)
	return http.StatusMovedPermanently, nil
}

func eveSSO(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	s.Values["state"] = state

	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	url := c.SSOAuthenticator.AuthorizeURL(state, true)
	http.Redirect(w, r, url, 302)
	return http.StatusMovedPermanently, nil
}

func eveSSOAnswer(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["state"] != state {

		return http.StatusInternalServerError, errors.New("Invalid State. It is possible that the session cookie is missing. Stop eating the cookies!")
	}

	tok, err := c.SSOAuthenticator.TokenExchange(c.HTTPClient, code)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed Token Exchange")
	}

	cli := c.SSOAuthenticator.GetClientFromToken(c.HTTPClient, tok)
	v, err := cli.Verify()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	s.Values["character"] = v
	s.Values["characterID"] = v.CharacterID
	s.Values["token"] = tok

	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/account", 302)
	return http.StatusMovedPermanently, nil
}

func eveCRESTToken(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	s.Values["TOKENstate"] = state

	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	url := c.TokenAuthenticator.AuthorizeURL(state, true)
	http.Redirect(w, r, url, 302)
	return http.StatusMovedPermanently, nil
}

func eveTokenAnswer(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["TOKENstate"] != state {
		return http.StatusInternalServerError, errors.New("Invalid State. It is possible that the session cookie is missing. Stop eating the cookies!")
	}

	tok, err := c.TokenAuthenticator.TokenExchange(c.HTTPClient, code)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed Token Exchange")
	}

	cli := c.TokenAuthenticator.GetClientFromToken(c.HTTPClient, tok)
	v, err := cli.Verify()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	characterID := s.Values["characterID"].(int64)
	err = models.AddCRESTToken(characterID, v.CharacterID, v.CharacterName, tok)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/account", 302)
	return http.StatusMovedPermanently, nil
}
