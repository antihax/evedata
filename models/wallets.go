package models

import "fmt"
import "time"

type WalletSummary struct {
	RefTypeID      int64          `db:"refTypeID" json:"refTypeID"`
	RefTypeName    string         `db:"refTypeName" json:"refTypeName"`
	Balance        float64        `db:"balance" json:"balance"`
	JournalEntries []JournalEntry `db:"journalEntries" json:"journalEntries,omitempty"`
}

type JournalEntry struct {
	RefId       int64     `db:"refId" json:"refId"`
	RefTypeID   int64     `db:"refTypeID" json:"refTypeID"`
	OwnerID1    int64     `db:"ownerID1" json:"ownerID1"`
	OwnerID2    int64     `db:"ownerID2" json:"ownerID2"`
	ArgID1      int64     `db:"argID1" json:"argID1"`
	ArgName1    string    `db:"argName1" json:"argName1"`
	Amount      float64   `db:"amount" json:"amount"`
	Reason      string    `db:"reason" json:"reason"`
	TaxAmount   float64   `db:"taxAmount" json:"taxAmount"`
	Date        time.Time `db:"date" json:"date"`
	CharacterID int64     `db:"characterID" json:"characterID"`
}

func GetWalletSummary(characterID int64, filterCharacterID int64) ([]WalletSummary, error) {

	filter := ""

	if filterCharacterID == 0 {
		filter = "IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND scopes LIKE '%wallet%')"
	} else {
		// False AST, forced int64.
		filter = fmt.Sprintf("IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND scopes LIKE '%%wallet%%' AND tokenCharacterID=%d)", filterCharacterID)
	}

	walletSummary := []WalletSummary{}
	if err := database.Select(&walletSummary, `
		SELECT T.refTypeID, refTypeName, SUM(amount) AS balance FROM evedata.walletJournal J
			INNER JOIN evedata.walletJournalRefType T ON J.refTypeID = T.refTypeID
			WHERE characterID `+filter+`
			GROUP BY refTypeID;
	`, characterID); err != nil {
		return nil, err
	}

	count := 0
	errc := make(chan error)
	limit := make(chan bool, 10)
	for index, _ := range walletSummary {
		count++
		go getJournalEntries(characterID, filterCharacterID, walletSummary[index].RefTypeID, &walletSummary[index].JournalEntries, errc, limit)
	}

	for i := 0; i < count; i++ {
		err := <-errc
		if err != nil {
			return nil, err // Something went wrong, break out.
		}
	}

	return walletSummary, nil
}

func getJournalEntries(characterID int64, filterCharacterID int64, refTypeID int64, entries *[]JournalEntry, errc chan error, limit chan bool) {
	limit <- true
	defer func() { <-limit }()

	filter := ""

	if filterCharacterID == 0 {
		filter = "IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND scopes LIKE '%wallet%')"
	} else {
		// False AST, forced int64.
		filter = fmt.Sprintf("IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND scopes LIKE '%%wallet%%' AND tokenCharacterID=%d)", filterCharacterID)
	}

	if err := database.Select(entries, `
		SELECT refID, refTypeID, ownerID1, ownerID2, argID1, argName1, amount, 
		reason, taxAmount, date 
		FROM evedata.walletJournal
		WHERE characterID `+filter+`
		AND refTypeID = ?
		ORDER BY date DESC;
	`, characterID, refTypeID); err != nil {
		errc <- err
		return
	}
}
