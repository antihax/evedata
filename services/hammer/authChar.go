package hammer

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterAuthOwner", characterAuthOwner)
	registerConsumer("corporationContacts", corporationContacts)
	registerConsumer("allianceContacts", allianceContacts)
}

func characterAuthOwner(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	roles, _, err := s.esi.ESI.CharacterApi.GetCharactersCharacterIdRoles(ctx, tokenCharacterID, nil)
	if err != nil {
		s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterRoles{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Roles:            roles,
	}, "characterAuthOwner")
	if err != nil {
		log.Println(err)
		return
	}
}

func allianceContacts(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))
	allianceID := int32(parameters[2].(int))

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	contacts, r, err := s.esi.ESI.ContactsApi.GetAlliancesAllianceIdContacts(ctx, allianceID, nil)
	if err != nil {
		return
	}

	pagesInt, err := strconv.Atoi(r.Header.Get("x-pages"))
	if err == nil {

		pages := int32(pagesInt)

		// Make a channel for reading contacts from the other pages
		conch := make(chan []esi.GetAlliancesAllianceIdContacts200Ok, 100)

		// Concurrently pull the remaining contact pages
		wg := sync.WaitGroup{}
		for pages != 1 {
			wg.Add(1)
			go func(page int32) {
				defer wg.Done()

				contacts, _, err := s.esi.ESI.ContactsApi.GetAlliancesAllianceIdContacts(ctx, allianceID,
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
	}

	// Send out the result
	err = s.QueueResult(&datapackages.AllianceContacts{
		AllianceID: allianceID,
		Contacts:   contacts,
	}, "allianceContacts")
	if err != nil {
		log.Println(err)
		return
	}
}

func corporationContacts(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))
	corporationID := int32(parameters[2].(int))

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	contacts, r, err := s.esi.ESI.ContactsApi.GetCorporationsCorporationIdContacts(ctx, corporationID, nil)
	if err != nil {
		return
	}

	pagesInt, err := strconv.Atoi(r.Header.Get("x-pages"))
	if err == nil {

		pages := int32(pagesInt)

		// Make a channel for reading contacts from the other pages
		conch := make(chan []esi.GetCorporationsCorporationIdContacts200Ok, 100)

		// Concurrently pull the remaining contact pages
		wg := sync.WaitGroup{}
		for pages != 1 {
			wg.Add(1)
			go func(page int32) {
				defer wg.Done()

				contacts, _, err := s.esi.ESI.ContactsApi.GetCorporationsCorporationIdContacts(ctx, corporationID,
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
	}
	// Send out the result
	err = s.QueueResult(&datapackages.CorporationContacts{
		CorporationID: corporationID,
		Contacts:      contacts,
	}, "corporationContacts")
	if err != nil {
		log.Println(err)
		return
	}
}
