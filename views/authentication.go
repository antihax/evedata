package views

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"evedata/eveapi"
	"evedata/models"
	"evedata/server"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute(evedata.Route{"logout", "GET", "/logout", logout})

	evedata.AddRoute(evedata.Route{"eveAuth", "GET", "/eveAuth", eveSSO})
	evedata.AddRoute(evedata.Route{"eveSSOAnswer", "GET", "/eveSSOAnswer", eveSSOAnswer})

	evedata.AddRoute(evedata.Route{"eveTokenAuth", "GET", "/eveTokenAuth", eveCRESTToken})
	evedata.AddRoute(evedata.Route{"eveTokenAnswer", "GET", "/eveTokenAnswer", eveTokenAnswer})
}

func logout(c *evedata.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	s.Options.MaxAge = -1
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/", 302)
	return http.StatusMovedPermanently, nil
}

func eveSSO(c *evedata.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	s.Values["state"] = state
	log.Printf("created state %v != %v\n", state, s.Values["state"])
	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	url := c.SSOAuthenticator.AuthorizeURL(state, true)
	http.Redirect(w, r, url, 302)
	return http.StatusMovedPermanently, nil
}

func eveSSOAnswer(c *evedata.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["state"] != state {
		log.Printf("confim state %v != %v\n", state, s.Values["state"])
		return http.StatusInternalServerError, errors.New("Invalid State. It is possible that the session cookie is missing. Stop eating the cookies!")
	}

	tok, err := c.SSOAuthenticator.TokenExchange(c.HTTPClient, code)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed Token Exchange")
	}
	cli := eveapi.NewAuthenticatedClient(c.HTTPClient, tok)
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

func eveCRESTToken(c *evedata.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
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

func eveTokenAnswer(c *evedata.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["TOKENstate"] != state {
		return http.StatusInternalServerError, errors.New("Invalid State. It is possible that the session cookie is missing. Stop eating the cookies!")
	}

	tok, err := c.TokenAuthenticator.TokenExchange(c.HTTPClient, code)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed Token Exchange")
	}

	cli := eveapi.NewAuthenticatedClient(c.HTTPClient, tok)
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
