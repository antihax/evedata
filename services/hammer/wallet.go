package hammer

import (
	"context"
	"log"

	"encoding/gob"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterWalletTransactions", characterWalletTransactionConsumer)
	registerConsumer("characterWalletJournal", characterWalletJournalConsumer)
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

func characterWalletJournalConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := parameters[0].(int32)
	tokenCharacterID := parameters[1].(int32)

	tokenSource, err := s.tokenStore.GetTokenSource(characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	journal, err := s.esi.EVEAPI.CharacterWalletJournalXML(tokenSource, int64(characterID), int64(tokenCharacterID))
	if err != nil {
		log.Println(err)
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterJournal{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Journal:          *journal,
	}, "characterWalletJournal")
	if err != nil {
		log.Println(err)
		return
	}
}
