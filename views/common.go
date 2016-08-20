package views

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type page struct {
	Title       string
	CharacterID int64
}

func newPage(s *sessions.Session, r *http.Request, title string) *page {
	p := &page{Title: title}

	characterID := s.Values["characterID"]

	if characterID != nil {
		p.CharacterID = characterID.(int64)
	} else {
		p.CharacterID = 0
	}

	return p
}
