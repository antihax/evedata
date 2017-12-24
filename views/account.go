package views

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/antihax/goesi"

	"github.com/antihax/evedata/evedata"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("account", "GET", "/account", accountPage)

	evedata.AddAuthRoute("account", "GET", "/X/accountInfo", accountInfo)
	evedata.AddAuthRoute("account", "POST", "/X/cursorChar", cursorChar)

	evedata.AddAuthRoute("crestTokens", "GET", "/U/crestTokens", apiGetCRESTTokens)
	evedata.AddAuthRoute("crestTokens", "DELETE", "/U/crestTokens", apiDeleteCRESTToken)

}

func accountPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	p := newPage(r, "Account Information")
	templates.Templates = template.Must(template.ParseFiles("templates/account.html", templates.LayoutPath))

	p["ScopeGroups"] = models.GetCharacterScopeGroups()

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		httpErr(w, err)
		return
	}
}

func accountInfo(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)

	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID"), http.StatusUnauthorized)
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
		if err := updateAccountInfo(s, int32(characterID), char.CharacterName); err != nil {
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
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

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
	if err = updateAccountInfo(s, characterID, char.CharacterName); err != nil {
		httpErr(w, err)
		return
	}

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiGetCRESTTokens(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for crest token"), http.StatusUnauthorized)
		return
	}

	v, err := models.GetCRESTTokens(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Change scopes to groups
	for i := range v {
		v[i].Scopes = models.GetCharacterGroupsByScopesString(v[i].Scopes)
	}

	json.NewEncoder(w).Encode(v)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}

func apiDeleteCRESTToken(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())
	g := evedata.GlobalsFromContext(r.Context())
	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID to delete"), http.StatusUnauthorized)
		return
	}

	cid, err := strconv.ParseInt(r.FormValue("tokenCharacterID"), 10, 64)
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		return
	}

	if err := models.DeleteCRESTToken(characterID, int32(cid)); err != nil {
		httpErrCode(w, err, http.StatusConflict)
		return
	}

	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find verify response to delete"), http.StatusForbidden)
		return
	}

	if err = updateAccountInfo(s, characterID, char.CharacterName); err != nil {
		httpErr(w, err)
		return
	}

	key := fmt.Sprintf("EVEDATA_TOKENSTORE_%d_%d", characterID, cid)
	red := g.Cache.Get()
	red.Do("DEL", key)
	red.Close()

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}
