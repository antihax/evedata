package hammer

import (
	"context"
	"log"
	"strconv"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
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

	var page int32 = 1
	assets := []esi.GetCharactersCharacterIdAssets200Ok{}
	for {
		a, r, err := s.esi.ESI.AssetsApi.GetCharactersCharacterIdAssets(ctx, tokenCharacterID,
			&esi.GetCharactersCharacterIdAssetsOpts{
				Page: optional.NewInt32(page),
			})
		if err != nil {
			s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
			log.Println(err)
			return
		}

		assets = append(assets, a...)

		xpagesS := r.Header.Get("x-pages")
		xpages, _ := strconv.Atoi(xpagesS)
		if int32(xpages) == page || len(a) == 0 {
			break
		}
		page++
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
