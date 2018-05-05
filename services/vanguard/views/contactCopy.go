package views

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/antihax/evedata/services/vanguard/models"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
)

func init() {
	vanguard.AddRoute("ContactCopy", "GET", "/contactCopy", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w,
			"contactCopy.html",
			time.Hour*24*31,
			newPage(r, "Contact Copy"))
	})
	vanguard.AddAuthRoute("ContactCopy", "PUT", "/U/contactCopy", apiPutContactCopy)
}

func apiPutContactCopy(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())
	g := vanguard.GlobalsFromContext(r.Context())
	char, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	source, err := strconv.Atoi(r.FormValue("source"))
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		log.Println(err)
		return
	}

	destination, err := strconv.Atoi(r.FormValue("destination"))
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		log.Println(err)
		return
	}

	if source == destination && r.FormValue("type") == "character" {
		httpErrCode(w, errors.New("Source and Destination cannot be the same"), http.StatusNotFound)
		return
	}

	tokenSource, err := g.TokenStore.GetTokenSource(char.CharacterID, int32(source))
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		log.Println(err)
		return
	}
	authSource := context.WithValue(context.Background(), goesi.ContextOAuth2, tokenSource)

	tokenDest, err := g.TokenStore.GetTokenSource(char.CharacterID, int32(destination))
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		log.Println(err)
		return
	}
	authDest := context.WithValue(context.Background(), goesi.ContextOAuth2, tokenDest)

	character, err := models.GetCRESTToken(char.CharacterID, char.CharacterOwnerHash, int32(source))
	if err != nil {
		httpErrCode(w, err, http.StatusNotFound)
		log.Println(err)
		return
	}

	var c []esi.GetCharactersCharacterIdContacts200Ok
	switch r.FormValue("type") {
	case "character":
		c, err = getCharacterContacts(g, authSource, int32(source))
		if err != nil {
			httpErrCode(w, err, http.StatusBadRequest)
			log.Println(err)
			return
		}

	case "corporation":
		c, err = getCorporationContacts(g, authSource, character.AllianceID)
		if err != nil {
			httpErrCode(w, err, http.StatusBadRequest)
			log.Println(err)
			return
		}

	case "alliance":
		c, err = getAllianceContacts(g, authSource, character.AllianceID)
		if err != nil {
			httpErrCode(w, err, http.StatusBadRequest)
			log.Println(err)
			return
		}

	default:
		httpErrCode(w, errors.New("Wrong Type!"), http.StatusBadRequest)
		return
	}

	toDelete, err := getCharacterContacts(g, authDest, int32(destination))
	if err != nil {
		httpErrCode(w, err, http.StatusBadRequest)
		log.Println(err)
		return
	}

	if err = eraseContacts(g, authDest, int32(destination), getContactIDs(toDelete)); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var p10, p5, neut, n5, n10 []int32
	for _, id := range c {
		switch id.Standing {
		case 10:
			p10 = append(p10, id.ContactId)
		case 5:
			p5 = append(p5, id.ContactId)
		case -10:
			n10 = append(n10, id.ContactId)
		case -5:
			n5 = append(n5, id.ContactId)
		default:
			neut = append(neut, id.ContactId)
		}
	}

	if err = addContacts(g, authDest, int32(destination), p10, 10); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err = addContacts(g, authDest, int32(destination), p5, 5); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err = addContacts(g, authDest, int32(destination), neut, 0); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err = addContacts(g, authDest, int32(destination), n10, -10); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err = addContacts(g, authDest, int32(destination), n5, -5); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func getContactIDs(contacts []esi.GetCharactersCharacterIdContacts200Ok) []int32 {
	ids := make([]int32, len(contacts))
	for i, c := range contacts {
		ids[i] = c.ContactId
	}
	return ids
}

func eraseContacts(g *vanguard.Vanguard, auth context.Context, characterID int32, erase []int32) error {
	// Erase contacts which have no wars.
	if len(erase) > 0 {
		for start := 0; start < len(erase); start = start + 20 {
			end := min(start+20, len(erase))
			if len(erase[start:end]) == 0 {
				break
			}
			if _, err := g.ESI.ESI.ContactsApi.DeleteCharactersCharacterIdContacts(auth, characterID, erase[start:end], nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func addContacts(g *vanguard.Vanguard, auth context.Context, characterID int32, c []int32, standing float32) error {
	// Add contacts for active wars
	if len(c) > 0 {
		for start := 0; start < len(c); start = start + 100 {
			end := min(start+100, len(c))
			if len(c[start:end]) == 0 {
				break
			}
			if _, _, err := g.ESI.ESI.ContactsApi.PostCharactersCharacterIdContacts(auth, characterID, c[start:end], standing, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func getCharacterContacts(g *vanguard.Vanguard, auth context.Context, characterID int32) ([]esi.GetCharactersCharacterIdContacts200Ok, error) {
	contacts, r, err := g.ESI.ESI.ContactsApi.GetCharactersCharacterIdContacts(auth, characterID, nil)
	if err != nil {
		return nil, err
	}

	// Decode the page into int32. Return if this fails as there were no extra pages.
	pagesInt, err := strconv.Atoi(r.Header.Get("x-pages"))
	if err != nil {
		return contacts, nil
	}
	pages := int32(pagesInt)

	// Make a channel for reading contacts from the other pages
	conch := make(chan []esi.GetCharactersCharacterIdContacts200Ok, 100)

	// Concurrently pull the remaining contact pages
	wg := sync.WaitGroup{}
	for {
		if pages == 1 {
			break
		}
		wg.Add(1)
		go func(page int32) {
			defer wg.Done()

			contacts, _, err := g.ESI.ESI.ContactsApi.GetCharactersCharacterIdContacts(auth, characterID,
				&esi.GetCharactersCharacterIdContactsOpts{Page: optional.NewInt32(page)})
			if err != nil {
				return
			}
			conch <- contacts
		}(pages)
		pages--
	}

	// Wait for everything to complete and close the channel.
	wg.Wait()
	close(conch)

	// Combine all the results
	for c := range conch {
		contacts = append(contacts, c...)
	}

	return contacts, nil
}

func getCorporationContacts(g *vanguard.Vanguard, auth context.Context, corporationID int32) ([]esi.GetCharactersCharacterIdContacts200Ok, error) {
	contacts, r, err := g.ESI.ESI.ContactsApi.GetCorporationsCorporationIdContacts(auth, corporationID, nil)
	if err != nil {
		return nil, err
	}

	// Decode the page into int32. Return if this fails as there were no extra pages.
	pages := int32(1)
	pagesInt, err := strconv.Atoi(r.Header.Get("x-pages"))
	if err == nil {
		pages = int32(pagesInt)
	}

	// Make a channel for reading contacts from the other pages
	conch := make(chan []esi.GetCorporationsCorporationIdContacts200Ok, 100)

	// Concurrently pull the remaining contact pages
	wg := sync.WaitGroup{}
	for {
		if pages == 1 {
			break
		}
		wg.Add(1)
		go func(page int32) {
			defer wg.Done()

			contacts, _, err := g.ESI.ESI.ContactsApi.GetCorporationsCorporationIdContacts(auth, corporationID,
				&esi.GetCorporationsCorporationIdContactsOpts{Page: optional.NewInt32(page)})
			if err != nil {
				return
			}
			conch <- contacts
		}(pages)
		pages--
	}

	// Wait for everything to complete and close the channel.
	wg.Wait()
	close(conch)

	// Combine all the results
	for c := range conch {
		contacts = append(contacts, c...)
	}

	newContacts := make([]esi.GetCharactersCharacterIdContacts200Ok, len(contacts))
	for i, c := range contacts {
		newContacts[i].ContactId = c.ContactId
		newContacts[i].Standing = c.Standing
	}

	return newContacts, nil
}

func getAllianceContacts(g *vanguard.Vanguard, auth context.Context, allianceID int32) ([]esi.GetCharactersCharacterIdContacts200Ok, error) {
	contacts, r, err := g.ESI.ESI.ContactsApi.GetAlliancesAllianceIdContacts(auth, allianceID, nil)
	if err != nil {
		return nil, err
	}

	// Decode the page into int32. Return if this fails as there were no extra pages.
	pages := int32(1)
	pagesInt, err := strconv.Atoi(r.Header.Get("x-pages"))
	if err == nil {
		pages = int32(pagesInt)
	}

	// Make a channel for reading contacts from the other pages
	conch := make(chan []esi.GetAlliancesAllianceIdContacts200Ok, 100)

	// Concurrently pull the remaining contact pages
	wg := sync.WaitGroup{}
	for {
		if pages == 1 {
			break
		}
		wg.Add(1)
		go func(page int32) {
			defer wg.Done()

			contacts, _, err := g.ESI.ESI.ContactsApi.GetAlliancesAllianceIdContacts(auth, allianceID,
				&esi.GetAlliancesAllianceIdContactsOpts{Page: optional.NewInt32(page)})
			if err != nil {
				return
			}
			conch <- contacts
		}(pages)
		pages--
	}

	// Wait for everything to complete and close the channel.
	wg.Wait()
	close(conch)

	// Combine all the results
	for c := range conch {
		contacts = append(contacts, c...)
	}

	newContacts := make([]esi.GetCharactersCharacterIdContacts200Ok, len(contacts))
	for i, c := range contacts {
		newContacts[i].ContactId = c.ContactId
		newContacts[i].Standing = c.Standing
	}

	return newContacts, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
