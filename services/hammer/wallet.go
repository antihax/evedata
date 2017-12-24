package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
)

func init() {
	registerConsumer("characterWalletTransactions", characterWalletTransactionConsumer)
	registerConsumer("characterWalletJournal", characterWalletJournalConsumer)
}

func characterWalletTransactionConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))

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
	if len(transactions) == 0 {
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterWalletTransactions{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Transactions:     transactions,
	}, "characterWalletTransactions")
	if err != nil {
		log.Println(err)
		return
	}
}

func characterWalletJournalConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	journal, _, err := s.esi.ESI.WalletApi.GetCharactersCharacterIdWalletJournal(ctx, tokenCharacterID, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if len(journal) == 0 {
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterJournal{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Journal:          journal,
	}, "characterWalletJournal")
	if err != nil {
		log.Println(err)
		return
	}
}
