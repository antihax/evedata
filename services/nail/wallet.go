package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/models"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("characterWalletTransactions", spawnCharacterWalletTransactionConsumer)
}

func spawnCharacterWalletTransactionConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterWalletTransactionConsumer)))
}

func (s *Nail) characterWalletTransactionConsumer(message *nsq.Message) error {
	wallet := datapackages.CharacterWalletTransactions{}
	err := gobcoder.GobDecoder(message.Body, &wallet)
	if err != nil {
		log.Println(err)
		return err
	}

	var values []string

	for _, transaction := range wallet.Transcations {
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
