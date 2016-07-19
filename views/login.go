package views

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"evedata/eveapi"
	"evedata/server"
	"net/http"
)

func init() {
	evedata.AddRoute(evedata.Route{"eveAuth", "GET", "/eveAuth", eveSSO})
	evedata.AddRoute(evedata.Route{"YourMotherWasAHampsterAndYourFatherSmeltOfElderberries", "GET", "/YourMotherWasAHampsterAndYourFatherSmeltOfElderberries", eveSSOAnswer})
}

func eveSSO(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	session, err := c.Store.Get(r, "session")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	session.Values["state"] = state
	session.Save(r, w)

	url := c.SSOAuthenticator.AuthorizeURL(state, true)
	http.Redirect(w, r, url, 302)
	return http.StatusMovedPermanently, nil
}

func eveSSOAnswer(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	code := r.FormValue("code")

	session, err := c.Store.Get(r, "session")
	if err != nil {
		return http.StatusInternalServerError, errors.New("Unable to connect to session store")
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

	session.Values["character"] = v
	session.Values["token"] = tok
	err = session.Save(r, w)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(w, r, "/", 302)
	return http.StatusMovedPermanently, nil
}
