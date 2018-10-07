package views

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/services/vanguard"

	"github.com/antihax/goesi"
	"github.com/gorilla/sessions"
)

// OpenGraph for Facebook unfurl
type OpenGraph struct {
	Title       string
	Image       string
	Description string
}

func entityImage(entityID int64, entityType string, size int) string {
	switch entityType {
	case "character":
		return "https://imageserver.eveonline.com/character/" +
			strconv.FormatInt(entityID, 10) +
			"_" + strconv.Itoa(size) + ".jpg"
	case "corporation":
		return "https://imageserver.eveonline.com/corporation/" +
			strconv.FormatInt(entityID, 10) +
			"_" + strconv.Itoa(size) + ".png"
	case "alliance":
		return "https://imageserver.eveonline.com/alliance/" +
			strconv.FormatInt(entityID, 10) +
			"_" + strconv.Itoa(size) + ".png"
	case "type":
		return "https://imageserver.eveonline.com/type/" +
			strconv.FormatInt(entityID, 10) +
			"_" + strconv.Itoa(size) + ".png"
	case "render":
		return "https://imageserver.eveonline.com/render/" +
			strconv.FormatInt(entityID, 10) +
			"_" + strconv.Itoa(size) + ".png"
	}
	return ""
}

func newPage(r *http.Request, title string) map[string]interface{} {
	p := make(map[string]interface{})
	p["Title"] = title
	p["Header"] = ""
	return p
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
