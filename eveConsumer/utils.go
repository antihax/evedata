package eveConsumer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"

	"golang.org/x/oauth2"
)

func (c *EVEConsumer) getContacts(auth context.Context, characterID int32) ([]esi.GetCharactersCharacterIdContacts200Ok, error) {
	var contacts []esi.GetCharactersCharacterIdContacts200Ok

	for i := int32(1); ; i++ {
		con, r, err := c.ctx.ESI.ESI.ContactsApi.GetCharactersCharacterIdContacts(auth, characterID, map[string]interface{}{"page": i})
		if err != nil || r.StatusCode != 200 {
			return c.getContactsCREST(auth, characterID)
		}
		if len(con) == 0 {
			break
		}
		contacts = append(contacts, con...)

		// Dirty hack to fix CCP breaking my stuff :(
		xpagesS := r.Header.Get("X-Pages")
		xpages, _ := strconv.Atoi(xpagesS)
		if int32(xpages) == i {
			break
		}

	}
	return contacts, nil
}

func (c *EVEConsumer) getContactsCREST(auth context.Context, characterID int32) ([]esi.GetCharactersCharacterIdContacts200Ok, error) {
	var contacts []esi.GetCharactersCharacterIdContacts200Ok

	tokenSource, ok := auth.Value(goesi.ContextOAuth2).(oauth2.TokenSource)
	if ok {

		con, err := c.ctx.ESI.EVEAPI.ContactsV1(tokenSource, int64(characterID))
		if err != nil {
			return nil, err
		}

		for ; con != nil; con, err = con.NextPage() {
			for _, contact := range con.Items {
				contacts = append(contacts,
					esi.GetCharactersCharacterIdContacts200Ok{
						ContactId:   int32(contact.Contact.ID),
						Standing:    float32(contact.Standing),
						IsBlocked:   contact.Blocked,
						IsWatched:   contact.Watched,
						ContactType: contact.ContactType,
					})
			}
		}
	}

	return contacts, nil
}

func (c *EVEConsumer) deleteContacts(auth context.Context, characterID int32, contacts []int32) error {
	_, err := c.ctx.ESI.ESI.ContactsApi.DeleteCharactersCharacterIdContacts(auth, characterID, contacts, nil)
	if err != nil {
		return c.deleteContactsCREST(auth, characterID, contacts)
	}
	return nil
}

func (c *EVEConsumer) deleteContactsCREST(auth context.Context, characterID int32, contacts []int32) error {
	names, _, err := c.ctx.ESI.ESI.UniverseApi.PostUniverseNames(nil, contacts, nil)
	if err != nil {
		return err
	}

	tokenSource, ok := auth.Value(goesi.ContextOAuth2).(oauth2.TokenSource)
	if ok {
		for _, erase := range names {
			ref := fmt.Sprintf("https://crest-tq.eveonline.com/%ss/%d/", erase.Category, erase.Id)
			err := c.ctx.ESI.EVEAPI.ContactDeleteV1(tokenSource, int64(characterID), int64(erase.Id), ref)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *EVEConsumer) addContacts(auth context.Context, characterID int32, contacts []int32, standing float32) error {
	_, _, err := c.ctx.ESI.ESI.ContactsApi.PostCharactersCharacterIdContacts(auth, characterID, contacts, standing, nil)
	if err != nil {
		return c.updateContacts(auth, characterID, contacts, standing)
	}
	return nil
}

func (c *EVEConsumer) updateContacts(auth context.Context, characterID int32, contacts []int32, standing float32) error {
	_, err := c.ctx.ESI.ESI.ContactsApi.PutCharactersCharacterIdContacts(auth, characterID, contacts, standing, nil)
	if err != nil {
		return c.updateContacts(auth, characterID, contacts, standing)
	}
	return nil
}

func (c *EVEConsumer) updateContactsCREST(auth context.Context, characterID int32, contacts []int32, standing float32) error {
	names, _, err := c.ctx.ESI.ESI.UniverseApi.PostUniverseNames(nil, contacts, nil)
	if err != nil {
		return err
	}

	tokenSource, ok := auth.Value(goesi.ContextOAuth2).(oauth2.TokenSource)
	if ok {
		for _, update := range names {
			ref := fmt.Sprintf("https://crest-tq.eveonline.com/%ss/%d/", update.Category, update.Id)
			err := c.ctx.ESI.EVEAPI.ContactSetV1(tokenSource, int64(characterID), int64(update.Id), ref, float64(standing))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func tokenError(cid int64, tcid int64, r *http.Response, err error) {
	if r != nil {
		req, _ := httputil.DumpRequest(r.Request, true)
		res, _ := httputil.DumpResponse(r, true)
		e := models.SetTokenError(cid, tcid, r.StatusCode, r.Status, req, res)
		if e != nil {
			log.Println(e)
		}
	} else {
		e := models.SetTokenError(cid, tcid, 999, err.Error(), []byte{}, []byte{})
		if e != nil {
			log.Println(e)
		}
	}
}

func tokenSuccess(cid int64, tcid int64, code int, status string) {
	models.SetTokenError(cid, tcid, code, status, []byte{}, []byte{})
}
