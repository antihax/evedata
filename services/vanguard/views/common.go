package views

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/services/vanguard"

	"github.com/antihax/goesi"
	"github.com/gorilla/sessions"
)

func newPage(r *http.Request, title string) map[string]interface{} {
	p := make(map[string]interface{})
	p["Title"] = title
	return p
}

func setCache(w http.ResponseWriter, cacheTime int) {
	if cacheTime == 0 {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	} else {
		w.Header().Set("Cache-Control", "max-age:"+strconv.Itoa(cacheTime)+", public")
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Expires", time.Now().UTC().Add(time.Second*time.Duration(cacheTime)).Format(http.TimeFormat))
	}
}

// getCursorCharacterAuth takes a session and returns the auth context or error
func getCursorCharacterAuth(ctx *vanguard.Vanguard, s *sessions.Session) (context.Context, error) {
	accountInfo, ok := s.Values["accountInfo"].([]byte)
	if !ok {
		return nil, errors.New("Cannot access account info")
	}

	info := accountInformation{}
	if err := json.Unmarshal(accountInfo, &info); err != nil {
		return nil, err
	}

	token, err := ctx.TokenStore.GetTokenSource(info.CharacterID, info.Cursor.CursorCharacterID)
	if err != nil {
		return nil, err
	}

	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, token)
	return auth, err
}

// getAccountInformation takes a session and returns the account information or error
func getAccountInformation(ctx *vanguard.Vanguard, s *sessions.Session) (*accountInformation, error) {
	accountInfo, ok := s.Values["accountInfo"].([]byte)
	if !ok {
		return nil, errors.New("Cannot access account info")
	}

	info := accountInformation{}
	if err := json.Unmarshal(accountInfo, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
