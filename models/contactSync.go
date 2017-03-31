package models

import (
	"errors"
	"time"

	"github.com/guregu/null"
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

func (c *ContactSync) Error(e string) error {
	_, err := database.Exec(`UPDATE evedata.contactSyncs SET lastError = ? WHERE source = ?`,
		e, c.Source)
	return err
}

func (c *ContactSync) Updated(nextSync time.Time) error {
	_, err := database.Exec(`UPDATE evedata.contactSyncs SET nextSync = ? WHERE source = ?`,
		nextSync, c.Source)
	return err
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetContactSyncs(characterID int64) ([]ContactSync, error) {
	cc := []ContactSync{}
	if err := database.Select(&cc, `
		SELECT C.characterID, source, S.characterName AS sourceName, destination, D.characterName AS destinationName, nextSync
			FROM evedata.contactSyncs C
	        LEFT JOIN evedata.crestTokens D ON C.destination = D.tokenCharacterID
			LEFT JOIN evedata.crestTokens S ON C.source = S.tokenCharacterID
			WHERE C.characterID = ?;`, characterID); err != nil {

		return nil, err
	}
	return cc, nil
}

func AddContactSync(characterID int64, source int, destination int) error {

	if source == destination {
		return errors.New("Source and Destination cannot be the same.")
	}
	if _, err := database.Exec(`INSERT IGNORE INTO evedata.contactSyncs (characterID, source, destination)VALUES(?,?,?)`,
		characterID, source, destination); err != nil {

		return err
	}
	return nil
}

func DeleteContactSync(characterID int64, destination int) error {
	if _, err := database.Exec(`DELETE FROM evedata.contactSyncs WHERE characterID = ? AND destination = ? LIMIT 1`,
		characterID, destination); err != nil {

		return err
	}
	return nil
}

type ExpiredContactSync struct {
	Source       int64  `db:"source" json:"source"`
	Destinations string `db:"destinations" json:"destinations"`
	CharacterID  int64  `db:"characterID" json:"characterID"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetExpiredContactSyncs() ([]ExpiredContactSync, error) {
	ecc := []ExpiredContactSync{}
	if err := database.Select(&ecc, `
		SELECT S.characterID, source, group_concat(destination) AS destinations
			FROM evedata.contactSyncs S  
            INNER JOIN evedata.crestTokens T ON T.tokenCharacterID = destination
            WHERE lastStatus NOT LIKE "%400 Bad Request%"
		    GROUP BY source
            HAVING max(nextSync) < UTC_TIMESTAMP();`); err != nil {

		return nil, err
	}
	return ecc, nil
}
