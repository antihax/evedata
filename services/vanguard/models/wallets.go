package models

import (
	"time"

	"github.com/guregu/null"
)

type WalletSummary struct {
	RefTypeID      int64          `db:"refTypeID" json:"refTypeID"`
	RefTypeName    string         `db:"refTypeName" json:"refTypeName"`
	Balance        float64        `db:"balance" json:"balance"`
	JournalEntries []JournalEntry `db:"journalEntries" json:"journalEntries,omitempty"`
}

type JournalEntry struct {
	RefID         int64       `db:"refID" json:"refID"`
	RefTypeID     int64       `db:"refTypeID" json:"refTypeID"`
	OwnerID1      int64       `db:"ownerID1" json:"ownerID1"`
	OwnerName1    null.String `db:"ownerName1" json:"ownerName1"`
	OwnerID2      int64       `db:"ownerID2" json:"ownerID2"`
	OwnerName2    null.String `db:"ownerName2" json:"ownerName2"`
	ArgID1        int64       `db:"argID1" json:"argID1"`
	ArgName1      string      `db:"argName1" json:"argName1"`
	Amount        float64     `db:"amount" json:"amount"`
	Reason        string      `db:"reason" json:"reason"`
	TaxAmount     float64     `db:"taxAmount" json:"taxAmount"`
	Date          time.Time   `db:"date" json:"date"`
	CharacterID   int64       `db:"characterID" json:"characterID"`
	CharacterName string      `db:"characterName" json:"characterName"`
}

func GetWalletSummary(characterID int32, rangeI int64) ([]WalletSummary, error) {

	walletSummary := []WalletSummary{}
	if err := database.Select(&walletSummary, `
		SELECT T.refTypeID, refTypeName, SUM(amount) AS balance FROM evedata.walletJournal J
			INNER JOIN evedata.walletJournalRefType T ON J.refTypeID = T.refTypeID
			WHERE characterID IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ?)
			AND date > DATE_SUB(UTC_TIMESTAMP(), INTERVAL ? DAY)
			GROUP BY refTypeID;
	`, characterID, rangeI); err != nil {
		return nil, err
	}

	count := 0
	errc := make(chan error)
	limit := make(chan bool, 20)
	for index := range walletSummary {
		count++
		go getJournalEntries(characterID, rangeI, walletSummary[index].RefTypeID, &walletSummary[index].JournalEntries, errc, limit)
	}

	for i := 0; i < count; i++ {
		err := <-errc
		if err != nil {
			return nil, err // Something went wrong, break out.
		}
	}

	return walletSummary, nil
}

func getJournalEntries(characterID int32, rangeI int64, refTypeID int64, entries *[]JournalEntry, errc chan error, limit chan bool) {
	limit <- true
	defer func() { <-limit }()

	if err := database.Select(entries, `
		SELECT refID, refTypeID, 
			ownerID1, C1.name AS ownerName1,
			ownerID2, C2.name AS ownerName2,
			argID1, argName1, amount, coalesce(typeName, stationName, reason) AS reason, taxAmount, 
			date, T.characterID, characterName
		FROM evedata.walletJournal J
		INNER JOIN evedata.crestTokens T ON J.characterID = T.tokenCharacterID
		LEFT JOIN evedata.characters C1 ON J.ownerID1 = C1.characterID
		LEFT JOIN evedata.characters C2 ON J.ownerID2 = C2.characterID
		LEFT JOIN evedata.walletTransactions TR ON J.refID = TR.journalTransactionID
        LEFT JOIN invTypes TY ON TY.typeID = TR.typeID
        LEFT JOIN staStations ST ON J.argID1 = ST.stationID
		WHERE T.characterID IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ?)
		AND refTypeID = ?
		AND date > DATE_SUB(UTC_TIMESTAMP(), INTERVAL ? DAY)
		ORDER BY date DESC;
	`, characterID, refTypeID, rangeI); err != nil {
		errc <- err
		return
	}

	errc <- nil
}
