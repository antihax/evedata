package views

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"

	"github.com/antihax/eveapi"
	"github.com/antihax/evedata/appContext"

	"github.com/antihax/evedata/models"
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

// Obtain an authenticated client from a stored access/refresh token.
func getToken(ctx *appContext.AppContext, characterID int64, tokenCharacterID int64) (oauth2.TokenSource, error) {
	tok, err := models.GetCRESTToken(characterID, tokenCharacterID)
	if err != nil {
		return nil, err
	}

	token := &eveapi.CRESTToken{Expiry: tok.Expiry, AccessToken: tok.AccessToken, RefreshToken: tok.RefreshToken, TokenType: tok.TokenType}
	n, err := ctx.TokenAuthenticator.TokenSource(token)

	return n, err
}

// getCursorCharacterAuth takes a session and returns the auth context or error
func getCursorCharacterAuth(ctx *appContext.AppContext, s *sessions.Session) (context.Context, error) {

	accountInfo, ok := s.Values["accountInfo"].([]byte)
	if !ok {
		return nil, errors.New("Cannot access account info")
	}

	info := accountInformation{}
	if err := json.Unmarshal(accountInfo, &info); err != nil {
		return nil, err
	}

	token, err := getToken(ctx, info.CharacterID, info.Cursor.CursorCharacterID)
	if err != nil {
		return nil, err
	}

	auth := context.WithValue(context.TODO(), goesi.ContextOAuth2, token)
	return auth, err
}
