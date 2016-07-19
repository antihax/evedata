package views

import (
	"evedata/server"
	"net/http"
)

type Page struct {
	Title       string
	CharacterID int
}

func NewPage(c *evedata.AppContext, r *http.Request, title string) *Page {
	p := &Page{Title: title}
	session, _ := c.Store.Get(r, "session")
	p.CharacterID = session.Values["characterID"].(int)
	return p
}
