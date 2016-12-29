package models

import (
	"errors"
	"time"

	"github.com/antihax/evedata/null"
)

type ContactSync struct {
	Source          int64       `db:"source" json:"source"`
	SourceName      null.String `db:"sourceName" json:"sourceName"`
	Destination     int64       `db:"destination" json:"destination"`
	DestinationName null.String `db:"destinationName" json:"destinationName"`
	CharacterID     int64       `db:"characterID" json:"characterID"`
	LastError       null.String `db:"lastError" json:"lastError"`
	NextSync        time.Time   `db:"nextSync" json:"nextSync"`
}

func (c *ContactSync) Error(err string) {
	database.Exec(`UPDATE contactSyncs SET lastError = ? WHERE source = ?`,
		err, c.Source)
}

func (c *ContactSync) Updated(nextSync time.Time) {
	database.Exec(`UPDATE contactSyncs SET nextSync = ? WHERE source = ?`,
		nextSync, c.Source)
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetContactSyncs(characterID int64) ([]ContactSync, error) {
	cc := []ContactSync{}
	if err := database.Select(&cc, `
		SELECT C.characterID, source, S.characterName AS sourceName, destination, D.characterName AS destinationName, nextSync
			FROM contactSyncs C
	        LEFT JOIN crestTokens D ON C.destination = D.tokenCharacterID
			LEFT JOIN crestTokens S ON C.source = S.tokenCharacterID
			WHERE C.characterID = ?;`, characterID); err != nil {

		return nil, err
	}
	return cc, nil
}

func AddContactSync(characterID int64, source int, destination int) error {

	if source == destination {
		return errors.New("Source and Destination cannot be the same.")
	}
	if _, err := database.Exec(`INSERT INTO contactSyncs (characterID, source, destination)VALUES(?,?,?)`,
		characterID, source, destination); err != nil {

		return err
	}
	return nil
}

func DeleteContactSync(characterID int64, destination int) error {
	if _, err := database.Exec(`DELETE FROM contactSyncs WHERE characterID = ? AND destination = ? LIMIT 1`,
		characterID, destination); err != nil {

		return err
	}
	return nil
}
