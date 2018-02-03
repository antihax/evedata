package models

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"golang.org/x/oauth2"
)

// Obtain an authenticated client from a stored access/refresh token.
func GetCRESTToken(characterID int32, tokenCharacterID int32) (*CRESTToken, error) {
	tok := &CRESTToken{}
	if err := database.QueryRowx(
		`SELECT expiry, tokenType, accessToken, refreshToken, tokenCharacterID, characterID, characterName
			FROM evedata.crestTokens
			WHERE characterID = ? AND tokenCharacterID = ?
			LIMIT 1`,
		characterID, tokenCharacterID).StructScan(tok); err != nil {

		return nil, err
	}

	return tok, nil
}

type CRESTToken struct {
	Expiry           time.Time   `db:"expiry" json:"expiry,omitempty"`
	CharacterID      int32       `db:"characterID" json:"characterID,omitempty"`
	TokenType        string      `db:"tokenType" json:"tokenType,omitempty"`
	TokenCharacterID int32       `db:"tokenCharacterID" json:"tokenCharacterID,omitempty"`
	CharacterName    string      `db:"characterName" json:"characterName,omitempty"`
	LastCode         int64       `db:"lastCode" json:"lastCode,omitempty"`
	LastStatus       null.String `db:"lastStatus" json:"lastStatus,omitempty"`
	AccessToken      string      `db:"accessToken" json:"accessToken,omitempty"`
	RefreshToken     string      `db:"refreshToken" json:"refreshToken,omitempty"`
	Scopes           string      `db:"scopes" json:"scopes"`
	AuthCharacter    int         `db:"authCharacter" json:"authCharacter"`
}

// [BENCHMARK] TODO
func GetCharacterIDByName(character string) (int32, error) {
	var id int32
	if err := database.Get(&id, `
		SELECT characterID 
		FROM evedata.characters C
		WHERE C.name = ? LIMIT 1;`, character); err != nil && err != sql.ErrNoRows {
		return id, err
	}
	return id, nil
}

type CursorCharacter struct {
	CursorCharacterID   int32  `db:"cursorCharacterID" json:"cursorCharacterID"`
	CursorCharacterName string `db:"cursorCharacterName" json:"cursorCharacterName"`
}

// [BENCHMARK] TODO
func GetCursorCharacter(characterID int32) (CursorCharacter, error) {
	cursor := CursorCharacter{}

	if err := database.Get(&cursor, `
		SELECT cursorCharacterID, T.characterName AS cursorCharacterName
		FROM evedata.cursorCharacter C
		INNER JOIN evedata.crestTokens T ON C.cursorCharacterID = T.tokenCharacterID AND C.characterID = T.characterID
		WHERE C.characterID = ?;`, characterID); err != nil {
		return cursor, err
	}
	return cursor, nil
}

// [BENCHMARK] TODO
func SetCursorCharacter(characterID int32, cursorCharacterID int32) error {
	if _, err := database.Exec(`
	INSERT INTO evedata.cursorCharacter (characterID, cursorCharacterID)
		SELECT characterID, tokenCharacterID AS cursorCharacterID
		FROM evedata.crestTokens WHERE characterID = ? AND tokenCharacterID = ? LIMIT 1
	ON DUPLICATE KEY UPDATE cursorCharacterID = VALUES(cursorCharacterID)
		;`, characterID, cursorCharacterID); err != nil {
		return err
	}
	return nil
}

func SetTokenError(characterID int32, tokenCharacterID int32, code int, status string, req []byte, res []byte) error {
	if _, err := database.Exec(`
		UPDATE evedata.crestTokens SET lastCode = ?, lastStatus = ?, request = ?, response = ? 
		WHERE characterID = ? AND tokenCharacterID = ? `,
		code, status, req, res, characterID, tokenCharacterID); err != nil {
		return err
	}
	return nil
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCRESTTokens(characterID int32) ([]CRESTToken, error) {
	tokens := []CRESTToken{}
	if err := database.Select(&tokens, `
		SELECT characterID, tokenCharacterID, characterName, lastCode, lastStatus, scopes, authCharacter
		FROM evedata.crestTokens
		WHERE characterID = ?;`, characterID); err != nil {

		return nil, err
	}
	return tokens, nil
}

func AddCRESTToken(characterID int32, tokenCharacterID int32, characterName string, tok *oauth2.Token, scopes string) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.crestTokens	(characterID, tokenCharacterID, accessToken, refreshToken, expiry, tokenType, characterName, scopes, lastStatus)
			VALUES		(?,?,?,?,?,?,?,?,"Unused")
			ON DUPLICATE KEY UPDATE 
				accessToken 	= VALUES(accessToken),
				refreshToken 	= VALUES(refreshToken),
				expiry 			= VALUES(expiry),
				tokenType 		= VALUES(tokenType),
				scopes 			= VALUES(scopes),
				lastStatus		= "Unused",
				mailedError 	= 0`,
		characterID, tokenCharacterID, tok.AccessToken, tok.RefreshToken, tok.Expiry, tok.TokenType, characterName, scopes); err != nil {
		return err
	}

	return nil
}

func DeleteCRESTToken(characterID int32, tokenCharacterID int32) error {
	if _, err := database.Exec(`DELETE FROM evedata.crestTokens WHERE characterID = ? AND tokenCharacterID = ? LIMIT 1`,
		characterID, tokenCharacterID); err != nil {

		return err
	}
	return nil
}

func UpdateCharacter(characterID int32, name string, bloodlineID int32, ancestryID int32, corporationID int32, allianceID int32,
	race int32, gender string, securityStatus float32, cacheUntil time.Time) error {
	cacheUntil = time.Now().UTC().Add(time.Hour * 24 * 5)
	if _, err := database.Exec(`
		INSERT INTO evedata.characters (characterID,name,bloodlineID,ancestryID,corporationID,allianceID,race,gender,securityStatus,updated,cacheUntil)
			VALUES(?,?,?,?,?,?,evedata.raceByID(?),?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
			corporationID=VALUES(corporationID), gender=VALUES(gender), allianceID=VALUES(allianceID), securityStatus=VALUES(securityStatus), updated = UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, characterID, name, bloodlineID, ancestryID, corporationID, allianceID, race, gender, securityStatus, cacheUntil); err != nil {
		return err
	}
	return nil
}

func UpdateCorporationHistory(characterID int32, corporationID int32, recordID int32, startDate time.Time) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.corporationHistory (characterID,startDate,recordID,corporationID)
			VALUES(?,?,?,?) 
			ON DUPLICATE KEY UPDATE 
			startDate=VALUES(startDate)
	`, characterID, startDate, recordID, corporationID); err != nil {
		return err
	}
	return nil
}

type Character struct {
	CharacterID     int32       `db:"characterID" json:"characterID"`
	CharacterName   string      `db:"characterName" json:"characterName"`
	CorporationID   int32       `db:"corporationID" json:"corporationID"`
	CorporationName string      `db:"corporationName" json:"corporationName"`
	AllianceID      int32       `db:"allianceID" json:"allianceID"`
	AllianceName    null.String `db:"allianceName" json:"allianceName"`
	Race            string      `db:"race" json:"race"`
	SecurityStatus  float64     `db:"securityStatus" json:"securityStatus"`
}

// Obtain Character information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCharacter(id int32) (*Character, error) {
	ref := Character{}
	if err := database.QueryRowx(`
		SELECT 
			characterID,
			C.name AS characterName,
		    C.corporationID,
		    IFNULL(Co.name, "Unknown Name") AS corporationName,
		    C.allianceID,
		    Al.name AS allianceName,
		    race,
		    securityStatus
		
		FROM evedata.characters C
		LEFT OUTER JOIN evedata.corporations Co ON Co.corporationID = C.corporationID
		LEFT OUTER JOIN evedata.alliances Al ON Al.allianceID = C.allianceID
		WHERE characterID = ?
		LIMIT 1`, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

type CorporationHistory struct {
	CorporationID   int32     `db:"corporationID" json:"id"`
	CorporationName string    `db:"corporationName" json:"name"`
	StartDate       time.Time `db:"startDate" json:"startDate"`
	Type            string    `db:"type" json:"type"`
}

// Obtain Character information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCorporationHistory(id int32) ([]CorporationHistory, error) {
	ref := []CorporationHistory{}
	if err := database.Select(&ref, `
		SELECT 
			C.corporationID,
			C.name AS corporationName,
			startDate
		    
		FROM evedata.corporationHistory H
		INNER JOIN evedata.corporations C ON C.corporationID = H.corporationID
		WHERE H.characterID = ?
		ORDER BY startDate DESC
		`, id); err != nil {
		return nil, err
	}
	for i := range ref {
		ref[i].Type = "corporation"
	}
	return ref, nil
}

type Shares struct {
	CharacterID        int32  `db:"characterID" json:"characterID,omitempty"`
	TokenCharacterID   int32  `db:"tokenCharacterID" json:"tokenCharacterID,omitempty"`
	TokenCharacterName string `db:"tokenCharacterName" json:"tokenCharacterName,omitempty"`
	EntityID           int32  `db:"entityID" json:"id,omitempty"`
	EntityName         string `db:"entityName" json:"entityName,omitempty"`
	Type               string `db:"type" json:"type,omitempty"`
	Types              string `db:"types" json:"types,omitempty"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetShares(characterID int32) ([]Shares, error) {
	shares := []Shares{}
	if err := database.Select(&shares, `
		SELECT S.characterID, S.tokenCharacterID, characterName AS tokenCharacterName, entityID, types, IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS type
		FROM evedata.sharing S
		INNER JOIN evedata.crestTokens T ON T.tokenCharacterID = S.tokenCharacterID AND T.characterID = S.characterID
		LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
		WHERE S.characterID = ?;`, characterID); err != nil {
		return nil, err
	}
	return shares, nil
}

func AddShare(characterID, tokenCharacterID, entityID int32, types string) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.sharing	(characterID, tokenCharacterID, entityID, types)
			VALUES(?,?,?,?)
			ON DUPLICATE KEY UPDATE entityID = entityID, types = VALUES(types)`,
		characterID, tokenCharacterID, entityID, types); err != nil {
		return err
	}
	return nil
}

func DeleteShare(characterID, tokenCharacterID, entityID int32) error {
	if _, err := database.Exec(`DELETE FROM evedata.sharing WHERE characterID = ? AND tokenCharacterID=? AND entityID = ? LIMIT 1`,
		characterID, tokenCharacterID, entityID); err != nil {
		return err
	}
	return nil
}
