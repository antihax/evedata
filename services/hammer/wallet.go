package hammer

import (
	"context"
	"log"

	"encoding/gob"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
)

func init() {
	registerConsumer("characterWalletTransactions", walletTransactionConsumer)
	gob.Register(datapackages.CharacterWalletTransactions{})
	gob.Register(datapackages.CharacterJournal{})
}

func walletTransactionConsumer(s *Hammer, parameter interface{}) {
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

	b, err := gobcoder.GobEncoder(&datapackages.CharacterWalletTransactions{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Transcations:     transactions,
	})

	if err != nil {
		log.Println(err)
		return
	}

	err = s.nsq.Publish("characterWalletTransactions", b)
	if err != nil {
		log.Println(err)
		return
	}
}
