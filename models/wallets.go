package models

import "fmt"

type WalletSummary struct {
	RefTypeID   int64      `db:"refTypeID" json:"refTypeID"`
	RefTypeName string     `db:"refTypeName" json:"refTypeName"`
	Balance           float64 `db:"balance" json:"balance"`
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
			WHERE characterID ` + filter + `
			GROUP BY refTypeID;
	`, characterID); err != nil {
		return nil, err
	}
	return walletSummary, nil
}
