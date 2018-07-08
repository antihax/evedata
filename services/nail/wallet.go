package nail

import (
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/goesi"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/services/vanguard/models"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("characterOrders", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterOrdersConsumer)))
	})
	AddHandler("characterWalletTransactions", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterWalletTransactionConsumer)))
	})
	AddHandler("characterWalletJournal", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddHandler(s.wait(nsq.HandlerFunc(s.characterWalletJournalConsumer)))
	})
}

func (s *Nail) characterOrdersConsumer(message *nsq.Message) error {
	orders := datapackages.CharacterOrders{}
	err := gobcoder.GobDecoder(message.Body, &orders)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(orders.Orders) == 0 {
		return nil
	}

	err = s.doSQL("DELETE FROM evedata.orders WHERE characterID = ?;", orders.TokenCharacterID)
	if err != nil {
		log.Println(err)
		return err
	}

	// Dump all orders into the DB.
	sql := sq.Insert("evedata.orders").Columns(
		"orderid", "characterID", "duration", "isBuyOrder", "isCorporation", "escrow",
		"issued", "locationID", "minVolume", "price", "orderRange", "regionID", "typeID",
		"volumeRemain", "volumeTotal",
	)
	for _, g := range orders.Orders {
		sql = sql.Values(
			g.OrderId, orders.TokenCharacterID, g.Duration, boolToInt(g.IsBuyOrder), boolToInt(g.IsCorporation), g.Escrow,
			g.Issued, g.LocationId, g.MinVolume, g.Price, g.Range_, g.RegionId, g.TypeId,
			g.VolumeRemain, g.VolumeTotal,
		)
	}

	sqlq, args, err := sql.ToSql()
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.doSQL(sqlq+" ON DUPLICATE KEY UPDATE orderid = orderid", args...)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func boolToInt(b bool) int8 {
	if b {
		return 1
	}
	return 0
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
