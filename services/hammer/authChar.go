package hammer

import (
	"context"
	"log"

	"github.com/antihax/goesi/esi"

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

	var page int32 = 1
	contacts := []esi.GetAlliancesAllianceIdContacts200Ok{}

	for {
		c, _, err := s.esi.ESI.ContactsApi.GetAlliancesAllianceIdContacts(ctx, allianceID, map[string]interface{}{"page": page})
		if err != nil {
			log.Println(err)
			return
		} else if len(c) == 0 { // end of the pages
			break
		}

		contacts = append(contacts, c...)

		page++
	}
	// early out if there are no orders
	if len(contacts) == 0 {
		return
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

	var page int32 = 1
	contacts := []esi.GetCorporationsCorporationIdContacts200Ok{}

	for {
		c, _, err := s.esi.ESI.ContactsApi.GetCorporationsCorporationIdContacts(ctx, corporationID, map[string]interface{}{"page": page})
		if err != nil {
			log.Println(err)
			return
		} else if len(c) == 0 { // end of the pages
			break
		}

		contacts = append(contacts, c...)

		page++
	}
	// early out if there are no orders
	if len(contacts) == 0 {
		return
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
