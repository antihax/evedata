package views

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/antihax/goesi"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/gorilla/sessions"
)

func init() {
	vanguard.AddAuthRoute("logout", "GET", "/X/logout", logout)

	vanguard.AddAuthRoute("eveAuth", "GET", "/X/eveAuth", eveSSO)
	vanguard.AddAuthRoute("eveSSOAnswer", "GET", "/X/eveSSOAnswer", eveSSOAnswer)

	vanguard.AddAuthRoute("eveTokenAuth", "GET", "/X/eveTokenAuth", eveCRESTToken)
	vanguard.AddAuthRoute("eveTokenAnswer", "GET", "/X/eveTokenAnswer", eveTokenAnswer)
}

func logout(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())
	s.Options.MaxAge = -1
	err := s.Save(r, w)
	if err != nil {
		httpErr(w, err)
		return
	}

	http.Redirect(w, r, "/", 302)
	httpErrCode(w, nil, http.StatusMovedPermanently)
}

func eveSSO(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())
	c := vanguard.GlobalsFromContext(r.Context())

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	s.Values["state"] = state

	err := s.Save(r, w)
	if err != nil {
		httpErr(w, err)
		return
	}

	url := c.SSOAuthenticator.AuthorizeURL(state, true, nil)
	http.Redirect(w, r, url, 302)
	httpErrCode(w, nil, http.StatusMovedPermanently)
}

func eveSSOAnswer(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())
	c := vanguard.GlobalsFromContext(r.Context())

	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["state"] != state {
		httpErr(w, errors.New("state does not match. We likely could not read the session cookie. Please make sure cookies are enabled."))
		return
	}

	tok, err := c.SSOAuthenticator.TokenExchange(code)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	tokSrc, err := c.SSOAuthenticator.TokenSource(tok)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	v, err := c.SSOAuthenticator.Verify(tokSrc)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	s.Values["character"] = v
	s.Values["characterID"] = v.CharacterID
	s.Values["token"] = tok
	s.Values["state"] = ""

	if err = updateAccountInfo(s, v.CharacterID, v.CharacterName); err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	if err = s.Save(r, w); err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	http.Redirect(w, r, "/account", 302)
	httpErrCode(w, nil, http.StatusMovedPermanently)
}

type accountInformation struct {
	CharacterID   int32                  `json:"characterID"`
	CharacterName string                 `json:"characterName"`
	Characters    []models.CRESTToken    `json:"characters"`
	Cursor        models.CursorCharacter `json:"cursor"`
}

func updateAccountInfo(s *sessions.Session, characterID int32, characterName string) error {
	var err error
	a := accountInformation{}

	a.CharacterName = characterName
	a.CharacterID = characterID
	a.Characters, err = models.GetCRESTTokens(characterID)
	if err != nil {
		log.Println(err)
	}

	a.Cursor, err = models.GetCursorCharacter(characterID)
	if err != nil {
		log.Println(err)
	}
	b, err := json.Marshal(a)
	s.Values["accountInfo"] = b

	return err
}

func eveCRESTToken(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())
	c := vanguard.GlobalsFromContext(r.Context())

	var scopes []string

	// Get the scopeGroups
	scopeGroupsTxt := r.FormValue("scopeGroups")

	if scopeGroupsTxt != "" {
		// split into []string
		scopeGroups := strings.Split(scopeGroupsTxt, ",")

		// Validate the scopeGroups are actually real
		validate := models.GetCharacterScopeGroups()
		for _, group := range scopeGroups {
			if validate[group] == "" {
				httpErrCode(w, nil, http.StatusBadRequest)
				return
			}
		}
		// Get the associated scopes to the groups
		scopes = models.GetCharacterScopesByGroups(scopeGroups)
	}

	// Hack to allow no scopes
	scopes = append(scopes, "publicData")

	// Make a code to validate on the return
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Save the code to our session store to compare later
	s.Values["TOKENstate"] = state
	err := s.Save(r, w)
	if err != nil {
		httpErr(w, err)
		return
	}

	url := c.TokenAuthenticator.AuthorizeURL(state, true, scopes)

	http.Redirect(w, r, url, 302)
	httpErrCode(w, nil, http.StatusMovedPermanently)
}

func eveTokenAnswer(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := vanguard.SessionFromContext(r.Context())
	c := vanguard.GlobalsFromContext(r.Context())

	code := r.FormValue("code")
	state := r.FormValue("state")

	if s.Values["TOKENstate"] != state {
		httpErr(w, errors.New("State does not match."))
		return
	}

	tok, err := c.TokenAuthenticator.TokenExchange(code)
	if err != nil {
		httpErr(w, err)
		return
	}

	tokSrc, err := c.SSOAuthenticator.TokenSource(tok)
	if err != nil {
		httpErr(w, err)
		return
	}

	v, err := c.SSOAuthenticator.Verify(tokSrc)
	if err != nil {
		httpErr(w, err)
		return
	}

	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErr(w, errors.New("cannot find character in store"))
		return
	}

	err = models.AddCRESTToken(char.CharacterID, v.CharacterID, v.CharacterName, tok, v.Scopes)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Invalidate cache
	key := fmt.Sprintf("EVEDATA_TOKENSTORE_%d_%d", char.CharacterID, v.CharacterID)
	red := c.Cache.Get()
	defer red.Close()
	red.Do("DEL", key)

	if err = updateAccountInfo(s, char.CharacterID, char.CharacterName); err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	s.Values["TOKENstate"] = ""
	if err := s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}

	http.Redirect(w, r, "/account", 302)
	httpErrCode(w, nil, http.StatusMovedPermanently)
}
