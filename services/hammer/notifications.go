package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterNotifications", characterNotificationsConsumer)
}

func characterNotificationsConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	notifications, _, err := s.esi.ESI.CharacterApi.GetCharactersCharacterIdNotifications(ctx, tokenCharacterID, nil)
	if err != nil {
		s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
		log.Println(err)
		return
	}
	if len(notifications) == 0 {
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterNotifications{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Notifications:    notifications,
	}, "characterNotifications")
	if err != nil {
		log.Println(err)
		return
	}
}
