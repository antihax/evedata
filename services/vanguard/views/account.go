package views

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/antihax/goesi"
	"golang.org/x/oauth2"

	"github.com/antihax/evedata/services/conservator"
	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("account", "GET", "/account",
		func(w http.ResponseWriter, r *http.Request) {
			p := newPage(r, "Account Information")
			p["ScopeGroups"] = models.GetCharacterScopeGroups()
			renderTemplate(w, "account.html", time.Hour*24*31, p)
		})

	vanguard.AddAuthRoute("account", "GET", "/X/accountInfo", accountInfo)
	vanguard.AddAuthRoute("account", "POST", "/X/cursorChar", cursorChar)

	vanguard.AddAuthRoute("account", "GET", "/U/crestTokens", apiGetCRESTTokens)
	vanguard.AddAuthRoute("account", "DELETE", "/U/crestTokens", apiDeleteCRESTToken)

	vanguard.AddAuthRoute("account", "GET", "/U/integrationTokens", apiGetIntegrationTokens)
	vanguard.AddAuthRoute("account", "DELETE", "/U/integrationTokens", apiDeleteIntegrationToken)

	vanguard.AddAuthRoute("account", "POST", "/U/toggleAuth", apiToggleAuth)

	vanguard.AddAuthRoute("account", "GET", "/U/accessableIntegrations", apiAccessableIntegrations)
	vanguard.AddAuthRoute("account", "POST", "/U/joinIntegration", apiJoinIntegration)

	vanguard.AddAuthRoute("account", "POST", "/U/setMailPassword", apiSetMailPassword)

}

func apiToggleAuth(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())
	if s == nil {
		httpErrCode(w, errors.New("Cannot find session"), http.StatusInternalServerError)
		return
	}
	g := vanguard.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for toggle auth"), http.StatusUnauthorized)
		return
	}

	// Parse the characterID
	tokenCharacterID, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusForbidden)
		return
	}

	_, err = g.Db.Exec("UPDATE evedata.crestTokens SET authCharacter = ! authCharacter WHERE characterID = ? and tokenCharacterID = ?", characterID, tokenCharacterID)
	if err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		return
	}
}

func accountInfo(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		// Silently fail.
		return
	}

	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find verify response"), http.StatusForbidden)
		log.Printf("%+v\n", s.Values["character"])
		return
	}

	accountInfo, ok := s.Values["accountInfo"].([]byte)
	if !ok {
		if err := updateAccountInfo(s, int32(characterID), char.CharacterOwnerHash, char.CharacterName); err != nil {
			httpErr(w, err)
			return
		}

		if err := s.Save(r, w); err != nil {
			httpErr(w, err)
			return
		}
	}

	w.Write(accountInfo)
}

func cursorChar(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for cursor"), http.StatusUnauthorized)
		return
	}

	// Parse the cursorCharacterID
	cursorCharacterID, err := strconv.ParseInt(r.FormValue("cursorCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusForbidden)
		return
	}

	// Set our new cursor
	err = models.SetCursorCharacter(characterID, int32(cursorCharacterID))
	if err != nil {
		httpErr(w, err)
		return
	}

	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find verify for cursor"), http.StatusForbidden)
		return
	}

	// Update the account information in redis
	if err = updateAccountInfo(s, characterID, char.CharacterOwnerHash, char.CharacterName); err != nil {
		httpErr(w, err)
		return
	}

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiGetCRESTTokens(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find character response"), http.StatusForbidden)
		log.Printf("%+v\n", s.Values["character"])
		return
	}

	v, err := models.GetCRESTTokens(char.CharacterID, char.CharacterOwnerHash)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Change scopes to groups
	for i := range v {
		v[i].Scopes = models.GetCharacterGroupsByScopesString(v[i].Scopes)
	}

	renderJSON(w, v, 0)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiDeleteCRESTToken(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())
	g := vanguard.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find verify response to delete"), http.StatusForbidden)
		return
	}

	cid, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	// Revoke the token before we delete. Do not error out if this fails.
	if tok, err := models.GetCRESTToken(char.CharacterID, char.CharacterOwnerHash, int32(cid)); err != nil {
		log.Println(err)
	} else {
		err = g.TokenAuthenticator.TokenRevoke(tok.RefreshToken)
		if err != nil {
			log.Println(err)
		}
	}

	if err := models.DeleteCRESTToken(char.CharacterID, int32(cid)); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}

	if err = updateAccountInfo(s, char.CharacterID, char.CharacterOwnerHash, char.CharacterName); err != nil {
		httpErr(w, err)
		return
	}

	key := fmt.Sprintf("EVEDATA_TOKENSTORE_%d_%d", char.CharacterID, cid)
	red := g.Cache.Get()
	defer red.Close()
	red.Do("DEL", key)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiAccessableIntegrations(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for integration token"), http.StatusUnauthorized)
		return
	}

	v, err := models.GetAvailableIntegrations(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, 0)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiJoinIntegration(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())
	g := vanguard.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for integration token"), http.StatusUnauthorized)
		return
	}

	integrationID, err := strconv.ParseInt(r.FormValue("integrationID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	i, err := models.GetIntegrationsForCharacter(characterID, int32(integrationID))
	if err != nil {
		httpErr(w, err)
		return
	}

	token := &oauth2.Token{
		Expiry:       i.Expiry,
		AccessToken:  i.AccessToken,
		RefreshToken: i.RefreshToken,
		TokenType:    "Bearer",
	}

	// refresh the token if it expired
	if token.Expiry.After(time.Now()) {
		src, err := g.DiscordAuthenticator.TokenSource(token)
		if err != nil {
			httpErr(w, err)
			return
		}
		token, err = src.Token()
		if err != nil {
			httpErr(w, err)
			return
		}
	}

	if err := g.RPCall("Conservator.JoinUser", conservator.JoinUser{
		IntegrationID: i.IntegrationID,
		AccessToken:   token.AccessToken,
		UserID:        i.IntegrationUserID,
		CharacterName: i.CharacterName,
		CharacterID:   i.TokenCharacterID,
	}, &ok); err != nil || !ok {
		httpErr(w, err)
		return
	}

}

func apiGetIntegrationTokens(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for integration token"), http.StatusUnauthorized)
		return
	}

	v, err := models.GetIntegrationTokens(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, 0)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiDeleteIntegrationToken(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for integration token to delete"), http.StatusUnauthorized)
		return
	}

	if err := models.DeleteIntegrationToken(r.FormValue("type"), characterID, r.FormValue("userID")); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}
}

func verifyPassword(s string) bool {
	var upper, lower, number bool
	for _, s := range s {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
		case unicode.IsLower(s):
			lower = true
		}
	}
	if upper && lower && number && len(s) >= 12 {
		return true
	}
	return false
}

func apiSetMailPassword(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())
	if s == nil {
		httpErrCode(w, errors.New("could not find session"), http.StatusUnauthorized)
		return
	}

	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find verify response to change mail password"), http.StatusForbidden)
		return
	}

	tokenCharacterID, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 32)
	if err != nil {
		httpErrCode(w, errors.New("invalid tokenCharacterID"), http.StatusBadRequest)
		return
	}

	if !verifyPassword(r.FormValue("password")) {
		httpErrCode(w, errors.New("Password must be at least 12 characters with one uppercase, one lowercase, and one number"), http.StatusBadRequest)
		return
	}

	if err := models.SetMailPassword(char.CharacterID, int32(tokenCharacterID), char.CharacterOwnerHash, r.FormValue("password")); err != nil {
		httpErr(w, err)
		return
	}
}
