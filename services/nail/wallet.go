package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/goesi"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/services/vanguard/models"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("characterWalletTransactions", spawnCharacterWalletTransactionConsumer)
	AddHandler("characterWalletJournal", spawnCharacterWalletJournalConsumer)
}

func spawnCharacterWalletTransactionConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterWalletTransactionConsumer)))
}

func spawnCharacterWalletJournalConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterWalletJournalConsumer)))
}

func (s *Nail) characterWalletTransactionConsumer(message *nsq.Message) error {
	wallet := datapackages.CharacterWalletTransactions{}
	err := gobcoder.GobDecoder(message.Body, &wallet)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(wallet.Transactions) == 0 {
		return nil
	}
	var values []string

	for _, transaction := range wallet.Transactions {
		var isFor, order string
		if transaction.IsPersonal {
			isFor = "personal"
		} else {
			isFor = "corporation"
		}
		if transaction.IsBuy {
			order = "buy"
		} else {
			order = "sell"
		}

		values = append(values, fmt.Sprintf("(%d,%d,%d,%d,%f,%d,%d,%q,%q,%d,%q)",
			wallet.TokenCharacterID, transaction.TransactionId, transaction.Quantity, transaction.TypeId, transaction.UnitPrice,
			transaction.ClientId, transaction.LocationId, order,
			isFor, transaction.JournalRefId, transaction.Date.UTC().Format(models.SQLTimeFormat)))
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.walletTransactions
									(characterID, transactionID, quantity, typeID, price,
									clientID,  stationID, transactionType,
									transactionFor, journalTransactionID, transactionDateTime)
									VALUES %s ON DUPLICATE KEY UPDATE characterID=characterID;`, strings.Join(values, ",\n"))

	return s.doSQL(stmt)

}

func (s *Nail) characterWalletJournalConsumer(message *nsq.Message) error {
	journal := datapackages.CharacterJournal{}
	err := gobcoder.GobDecoder(message.Body, &journal)
	if err != nil {
		log.Println(err)
		return err
	}

	if len(journal.Journal) == 0 {
		return nil
	}

	var values []string

	for _, wallet := range journal.Journal {
		values = append(values, fmt.Sprintf("(%d,%d,%d,%d,%d,%d,%q,%f,%f,%q,%d,%f,%q)",
			journal.TokenCharacterID, wallet.Id, goesi.GetJournalRefID(wallet.RefType), wallet.FirstPartyId, wallet.SecondPartyId,
			wallet.ContextId, wallet.ContextIdType, wallet.Amount, wallet.Balance,
			wallet.Reason, wallet.TaxReceiverId, wallet.Tax, wallet.Date.UTC().Format(models.SQLTimeFormat)))
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.walletJournal
							(characterID, refID, refTypeID, ownerID1, ownerID2,
							argID1, argName1, amount, balance,
							reason, taxReceiverID, taxAmount, date)
							VALUES %s ON DUPLICATE KEY UPDATE characterID=characterID;`, strings.Join(values, ",\n"))

	return s.doSQL(stmt)

}
