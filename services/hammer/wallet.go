package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
)

func init() {
	registerConsumer("characterWalletTransactions", characterWalletTransactionConsumer)
	registerConsumer("characterWalletJournal", characterWalletJournalConsumer)
	registerConsumer("characterOrders", characterOrdersConsumer)
}

func characterOrdersConsumer(s *Hammer, parameter interface{}) {
	// dereference the parameters
	parameters := parameter.([]interface{})
	characterID := int32(parameters[0].(int))
	tokenCharacterID := int32(parameters[1].(int))

	ctx, err := s.GetTokenSourceContext(context.Background(), characterID, tokenCharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	orders, _, err := s.esi.ESI.MarketApi.GetCharactersCharacterIdOrders(ctx, tokenCharacterID, nil)
	if err != nil {
		s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
		log.Println(err)
		return
	}
	if len(orders) == 0 {
		return
	}

	// Send out the result
	err = s.QueueResult(&datapackages.CharacterOrders{
		CharacterID:      characterID,
		TokenCharacterID: tokenCharacterID,
		Orders:           orders,
	}, "characterOrders")
	if err != nil {
		log.Println(err)
		return
	}
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
		s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
		log.Println(err)
		return
	}
	if len(transactions) == 0 {
		return
	}

	last := lowestTransactionID(transactions)
	for {
		top, _, err := s.esi.ESI.WalletApi.GetCharactersCharacterIdWalletTransactions(ctx, tokenCharacterID,
			&esi.GetCharactersCharacterIdWalletTransactionsOpts{
				FromId: optional.NewInt64(last),
			})
		if err != nil {
			s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
			log.Println(err)
			return
		}
		if len(top) == 0 {
			break
		}

		transactions = append(transactions, top...)
		newlast := lowestTransactionID(top)
		if newlast == last {
			break
		}
		last = newlast
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

	page := int32(1)
	journal := []esi.GetCharactersCharacterIdWalletJournal200Ok{}
	for {
		top, _, err := s.esi.ESI.WalletApi.GetCharactersCharacterIdWalletJournal(ctx, tokenCharacterID,
			&esi.GetCharactersCharacterIdWalletJournalOpts{
				Page: optional.NewInt32(page),
			})
		if err != nil {
			s.tokenStore.CheckSSOError(characterID, tokenCharacterID, err)
			log.Println(err)
			return
		}
		if len(top) == 0 {
			break
		}
		page++
		journal = append(journal, top...)
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

func lowestTransactionID(j []esi.GetCharactersCharacterIdWalletTransactions200Ok) int64 {
	lowest := j[0].TransactionId
	for _, i := range j {
		if i.TransactionId < lowest {
			lowest = i.TransactionId
		}
	}
	return lowest
}
