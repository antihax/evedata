package hammer

import (
	"context"
	"log"

	"encoding/gob"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterWalletTransactions", characterWalletTransactionConsumer)
	gob.Register(datapackages.CharacterWalletTransactions{})
	gob.Register(datapackages.CharacterJournal{})
}

func characterWalletTransactionConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := parameters[0].(int32)
	tokenCharacterID := parameters[1].(int32)

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	transactions, _, err := s.esi.ESI.WalletApi.GetCharactersCharacterIdWalletTransactions(ctx, tokenCharacterID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterWalletTransactions{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Transcations:     transactions,
	}, "characterWalletTransactions")
	if err != nil {
		log.Println(err)
		return
	}
}
