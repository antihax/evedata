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
		var (
			argID   int64
			argName string
		)

		e := wallet.ExtraInfo

		if e.AllianceId != 0 {
			argID = int64(e.AllianceId)
			argName = "alliance"
		} else if e.CharacterId != 0 {
			argID = int64(e.CharacterId)
			argName = "character"
		} else if e.CorporationId != 0 {
			argID = int64(e.CorporationId)
			argName = "corporation"
		} else if e.ContractId != 0 {
			argID = int64(e.ContractId)
			argName = "contract"
		} else if e.DestroyedShipTypeId != 0 {
			argID = int64(e.DestroyedShipTypeId)
			argName = "destroyedship"
		} else if e.JobId != 0 {
			argID = int64(e.JobId)
			argName = "job"
		} else if e.LocationId != 0 {
			argID = int64(e.LocationId)
			argName = "location"
		} else if e.NpcId != 0 {
			argID = int64(e.NpcId)
			argName = e.NpcName
		} else if e.PlanetId != 0 {
			argID = int64(e.PlanetId)
			argName = "planet"
		} else if e.SystemId != 0 {
			argID = int64(e.SystemId)
			argName = "system"
		} else if e.TransactionId != 0 {
			argID = int64(e.TransactionId)
			argName = "transaction"
		}

		values = append(values, fmt.Sprintf("(%d,%d,%d,%d,%d,%d,%q,%f,%f,%q,%d,%f,%q)",
			journal.TokenCharacterID, wallet.RefId, goesi.GetJournalRefID(wallet.RefType), wallet.FirstPartyId, wallet.SecondPartyId,
			argID, argName, wallet.Amount, wallet.Balance,
			wallet.Reason, wallet.TaxReceiverId, wallet.Tax, wallet.Date.UTC().Format(models.SQLTimeFormat)))
	}

	stmt := fmt.Sprintf(`INSERT INTO evedata.walletJournal
							(characterID, refID, refTypeID, ownerID1, ownerID2,
							argID1, argName1, amount, balance,
							reason, taxReceiverID, taxAmount, date)
							VALUES %s ON DUPLICATE KEY UPDATE characterID=characterID;`, strings.Join(values, ",\n"))

	return s.doSQL(stmt)

}
