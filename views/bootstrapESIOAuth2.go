package views

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/server"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddAuthRoute("boostrap", "GET", "/X/boostrapEveAuth", boostrapEveSSO)
	evedata.AddAuthRoute("boostrap", "GET", "/X/boostrapEveSSOAnswer", boostrapEveSSOAnswer)
}

func boostrapEveSSO(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	s.Values["BOOTSTRAPstate"] = state

	err := s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	url := c.ESIBootstrapAuthenticator.AuthorizeURL(state, true, nil)
	http.Redirect(w, r, url, 302)
	return http.StatusMovedPermanently, nil
}

func boostrapEveSSOAnswer(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["BOOTSTRAPstate"] != state {

		return http.StatusInternalServerError, errors.New("Invalid State. It is possible that the session cookie is missing. Stop eating the cookies!")
	}

	tok, err := c.ESIBootstrapAuthenticator.TokenExchange(code)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tokSrc, err := c.ESIBootstrapAuthenticator.TokenSource(tok)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = c.EVE.Verify(tokSrc)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	s.Values["BOOTSTRAP"] = tok

	fmt.Fprintf(w, "%+v\n", tok)

	err = s.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
