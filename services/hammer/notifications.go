package hammer

import (
	"context"
	"log"

	"encoding/gob"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterNotifications", characterNotificationsConsumer)
	gob.Register(datapackages.CharacterNotifications{})
}

func characterNotificationsConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]int32)
	characterID := parameters[0]
	tokenCharacterID := parameters[1]

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	notifications, _, err := s.esi.ESI.CharacterApi.GetCharactersCharacterIdNotifications(ctx, tokenCharacterID, nil)
	if err != nil {
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
