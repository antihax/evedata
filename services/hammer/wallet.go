package hammer

import (
	"context"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/goesi/esi"
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
	log.Printf("journal got %d \n", len(journal))
	if len(journal) == 0 {
		return
	}

	last := lowestRefID(journal)
	for {
		top, _, err := s.esi.ESI.WalletApi.GetCharactersCharacterIdWalletJournal(ctx, tokenCharacterID,
			map[string]interface{}{"fromId": last})
		if err != nil {
			log.Println(err)
			return
		}
		if len(top) == 0 {
			break
		}
		log.Printf("journal got %d additional %d\n", len(top), last)
		journal = append(journal, top...)
		newlast := lowestRefID(top)
		if newlast == last {
			break
		}
		last = newlast
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

func lowestRefID(j []esi.GetCharactersCharacterIdWalletJournal200Ok) int64 {
	lowest := j[0].RefId
	for _, i := range j {
		if i.RefId < lowest {
			lowest = i.RefId
		}
	}
	return lowest
}
