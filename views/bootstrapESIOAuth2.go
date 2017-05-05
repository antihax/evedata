package views

import (
	"fmt"
	"net/http"

	"github.com/antihax/evedata/evedata"
)

// Bootstrap Token provides access to public information in EVE Online. This is not used by users.
func init() {
	evedata.AddAuthRoute("bootstrap", "GET", "/X/bootstrapEveAuth", bootstrapEveSSO)
	evedata.AddAuthRoute("bootstrap", "GET", "/X/bootstrapEveSSOAnswer", bootstrapEveSSOAnswer)
}

func bootstrapEveSSO(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)

	c := evedata.GlobalsFromContext(r.Context())

	url := c.ESIBootstrapAuthenticator.AuthorizeURL("NONE", true, nil)
	http.Redirect(w, r, url, 302)
	httpErrCode(w, http.StatusMovedPermanently)
}

func bootstrapEveSSOAnswer(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	c := evedata.GlobalsFromContext(r.Context())

	code := r.FormValue("code")

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

	fmt.Fprintf(w, "%+v\n", tok)
}
