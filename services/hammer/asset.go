package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterAssets", characterAssetsConsumer)
}

func characterAssetsConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	assets, _, err := s.esi.ESI.AssetsApi.GetCharactersCharacterIdAssets(ctx, tokenCharacterID, nil)
	if err != nil {
		s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
		log.Println(err)
		return
	}
	if len(assets) == 0 {
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterAssets{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Assets:           assets,
	}, "characterAssets")
	if err != nil {
		log.Println(err)
		return
	}
}
