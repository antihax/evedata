package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterAuthOwner", characterAuthOwner)

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
