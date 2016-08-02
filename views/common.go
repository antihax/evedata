package views

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Page struct {
	Title       string
	CharacterID int64
}

func NewPage(s *sessions.Session, r *http.Request, title string) *Page {
	p := &Page{Title: title}

	characterID := s.Values["characterID"]

	if characterID != nil {
		p.CharacterID = characterID.(int64)
	} else {
		p.CharacterID = 0
	}

	return p
}
