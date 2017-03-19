package models

import (
	"database/sql"
	"time"

	"github.com/antihax/goesi"
	"github.com/guregu/null"
)

// Obtain an authenticated client from a stored access/refresh token.
func GetCRESTToken(characterID int64, tokenCharacterID int64) (*CRESTToken, error) {
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
	CharacterID      int64       `db:"characterID" json:"characterID,omitempty"`
	TokenType        string      `db:"tokenType" json:"tokenType,omitempty"`
	TokenCharacterID int64       `db:"tokenCharacterID" json:"tokenCharacterID,omitempty"`
	CharacterName    string      `db:"characterName" json:"characterName,omitempty"`
	LastCode         int64       `db:"lastCode" json:"lastCode,omitempty"`
	LastStatus       null.String `db:"lastStatus" json:"lastStatus,omitempty"`
	AccessToken      string      `db:"accessToken" json:"accessToken,omitempty"`
	RefreshToken     string      `db:"refreshToken" json:"refreshToken,omitempty"`
	Scopes           string      `db:"scopes" json:"scopes,omitempty"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCRESTTokens(characterID int64) ([]CRESTToken, error) {
	tokens := []CRESTToken{}
	if err := database.Select(&tokens, `
		SELECT characterID, tokenCharacterID, characterName, lastCode, lastStatus, scopes
		FROM evedata.crestTokens
		WHERE characterID = ?;`, characterID); err != nil {

		return nil, err
	}
	return tokens, nil
}

// [BENCHMARK] TODO
func GetCharacterIDByName(character string) (int64, error) {
	var id int64
	if err := database.Get(&id, `
		SELECT characterID 
		FROM evedata.characters C
		WHERE C.name = ? LIMIT 1;`, character); err != nil && err != sql.ErrNoRows {
		return id, err
	}
	return id, nil
}

type CursorCharacter struct {
	CursorCharacterID   int64  `db:"cursorCharacterID" json:"cursorCharacterID"`
	CursorCharacterName string `db:"cursorCharacterName" json:"cursorCharacterName"`
}

// [BENCHMARK] TODO
func GetCursorCharacter(characterID int64) (CursorCharacter, error) {
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
func SetCursorCharacter(characterID int64, cursorCharacterID int64) error {
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

func SetTokenError(characterID int64, tokenCharacterID int64, code int, status string, req []byte, res []byte) error {
	if _, err := database.Exec(`
		UPDATE evedata.crestTokens SET lastCode = ?, lastStatus = ?, request = ?, response = ? 
		WHERE characterID = ? AND tokenCharacterID = ? `,
		code, status, req, res, characterID, tokenCharacterID); err != nil {
		return err
	}
	return nil
}

func AddCRESTToken(characterID int64, tokenCharacterID int64, characterName string, tok *goesi.CRESTToken, scopes string) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.crestTokens	(characterID, tokenCharacterID, accessToken, refreshToken, expiry, tokenType, characterName, scopes, lastStatus)
			VALUES		(?,?,?,?,?,?,?,?,"Unused")
			ON DUPLICATE KEY UPDATE 
				accessToken 	= VALUES(accessToken),
				refreshToken 	= VALUES(refreshToken),
				expiry 			= VALUES(expiry),
				tokenType 		= VALUES(tokenType),
				scopes 		= VALUES(scopes),
				lastStatus		= "Unused"`,
		characterID, tokenCharacterID, tok.AccessToken, tok.RefreshToken, tok.Expiry, tok.TokenType, characterName, scopes); err != nil {
		return err
	}

	return nil
}

func DeleteCRESTToken(characterID int64, tokenCharacterID int64) error {
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
	CharacterID     int64       `db:"characterID" json:"characterID"`
	CharacterName   string      `db:"characterName" json:"characterName"`
	CorporationID   int64       `db:"corporationID" json:"corporationID"`
	CorporationName string      `db:"corporationName" json:"corporationName"`
	AllianceID      int64       `db:"allianceID" json:"allianceID"`
	AllianceName    null.String `db:"allianceName" json:"allianceName"`
	Race            string      `db:"race" json:"race"`
	SecurityStatus  float64     `db:"securityStatus" json:"securityStatus"`
}

// Obtain Character information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCharacter(id int64) (*Character, error) {
	ref := Character{}
	if err := database.QueryRowx(`
		SELECT 
			characterID,
			C.name AS characterName,
		    C.corporationID,
		    Co.name AS corporationName,
		    C.allianceID,
		    Al.name AS allianceName,
		    race,
		    securityStatus
		
		FROM evedata.characters C
		INNER JOIN evedata.corporations Co ON Co.corporationID = C.corporationID
		LEFT OUTER JOIN evedata.alliances Al ON Al.allianceID = C.allianceID
		WHERE characterID = ?
		LIMIT 1`, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}
