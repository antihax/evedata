package views

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func newPage(s *sessions.Session, r *http.Request, title string) map[string]interface{} {
	p := make(map[string]interface{})

	p["Title"] = title
	characterID := s.Values["characterID"]

	if characterID != nil {
		p["CharacterID"] = characterID.(int64)
	} else {
		p["CharacterID"] = 0
	}

	return p
}
