package hammer

import (
	"context"
	"encoding/gob"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterAssets", characterAssetsConsumer)
	gob.Register(datapackages.CharacterAssets{})
}

func characterAssetsConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]int32)
	characterID := parameters[0]
	tokenCharacterID := parameters[1]

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	assets, _, err := s.esi.ESI.AssetsApi.GetCharactersCharacterIdAssets(ctx, tokenCharacterID, nil)
	if err != nil {
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
