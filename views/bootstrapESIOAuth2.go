package views

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/antihax/evedata/evedata"
)

func init() {
	evedata.AddAuthRoute("boostrap", "GET", "/X/boostrapEveAuth", boostrapEveSSO)
	evedata.AddAuthRoute("boostrap", "GET", "/X/boostrapEveSSOAnswer", boostrapEveSSOAnswer)
}

func boostrapEveSSO(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())
	c := evedata.GlobalsFromContext(r.Context())

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	s.Values["BOOTSTRAPstate"] = state

	err := s.Save(r, w)
	if err != nil {
		httpErr(w, err)
		return
	}

	url := c.ESIBootstrapAuthenticator.AuthorizeURL(state, true, nil)
	http.Redirect(w, r, url, 302)
	httpErrCode(w, http.StatusMovedPermanently)
}

func boostrapEveSSOAnswer(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())
	c := evedata.GlobalsFromContext(r.Context())

	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["BOOTSTRAPstate"] != state {
		httpErr(w, errors.New("Invalid State. It is possible that the session cookie is missing. Stop eating the cookies!"))
		return
	}

	tok, err := c.ESIBootstrapAuthenticator.TokenExchange(code)
	if err != nil {
		httpErr(w, err)
		return
	}

	tokSrc, err := c.ESIBootstrapAuthenticator.TokenSource(tok)
	if err != nil {
		httpErr(w, err)
		return
	}

	_, err = c.SSOAuthenticator.Verify(tokSrc)
	if err != nil {
		httpErr(w, err)
		return
	}

	if err != nil {
		httpErr(w, err)
		return
	}

	s.Values["BOOTSTRAP"] = tok

	fmt.Fprintf(w, "%+v\n", tok)

	err = s.Save(r, w)
	if err != nil {
		httpErr(w, err)
		return
	}
}
