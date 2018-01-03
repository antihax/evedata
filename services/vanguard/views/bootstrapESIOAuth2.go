package views

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/antihax/evedata/services/vanguard"
)

// Bootstrap Token provides access to public information in EVE Online. This is not used by users.
func init() {
	vanguard.AddAuthRoute("bootstrap", "GET", "/X/bootstrapEveAuth", bootstrapEveSSO)
	vanguard.AddAuthRoute("bootstrap", "GET", "/X/bootstrapEveSSOAnswer", bootstrapEveSSOAnswer)
}

func bootstrapEveSSO(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)

	c := vanguard.GlobalsFromContext(r.Context())
	bootstrapScopes := strings.Split("esi-calendar.respond_calendar_events.v1 esi-calendar.read_calendar_events.v1 esi-mail.organize_mail.v1 esi-mail.read_mail.v1 esi-mail.send_mail.v1 esi-wallet.read_character_wallet.v1 esi-wallet.read_corporation_wallet.v1 esi-search.search_structures.v1 esi-universe.read_structures.v1 esi-corporations.read_corporation_membership.v1 esi-markets.structure_markets.v1 esi-characters.read_chat_channels.v1 esi-corporations.track_members.v1 esi-wallet.read_corporation_wallets.v1 esi-corporations.read_divisions.v1 esi-assets.read_corporation_assets.v1", " ")

	url := c.PrivateAuthenticator.AuthorizeURL("NONE", true, bootstrapScopes)
	http.Redirect(w, r, url, 302)
	httpErrCode(w, nil, http.StatusMovedPermanently)
}

func bootstrapEveSSOAnswer(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	c := vanguard.GlobalsFromContext(r.Context())

	code := r.FormValue("code")

	tok, err := c.PrivateAuthenticator.TokenExchange(code)
	if err != nil {
		httpErr(w, err)
		return
	}

	tokSrc, err := c.PrivateAuthenticator.TokenSource(tok)
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
