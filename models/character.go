package models

import (
	"evedata/eveapi"
	"evedata/null"
	"time"
)

type ApiKey struct {
	KeyID       int64       `db:"keyID" json:"keyID"`
	CharacterID int64       `db:"characterID" json:"characterID"`
	NextCheck   time.Time   `db:"nextCheck" json:"nextCheck"`
	LastCode    int64       `db:"lastCode" json:"lastCode"`
	LastError   null.String `db:"lastError" json:"lastError"`
	AccessMask  int64       `db:"accessMask" json:"accessMask"`
	Type        null.String `db:"type" json:"type"`
}

type CRESTToken struct {
	Expiry           time.Time   `db:"expiry" json:"expiry"`
	CharacterID      int64       `db:"characterID" json:"characterID"`
	TokenType        string      `db:"tokenType" json:"tokenType"`
	TokenCharacterID int64       `db:"tokenCharacterID" json:"tokenCharacterID"`
	CharacterName    string      `db:"characterName" json:"characterName"`
	LastCode         int64       `db:"lastCode" json:"lastCode"`
	LastStatus       null.String `db:"lastStatus" json:"lastStatus"`
	AccessToken      string      `db:"accessToken" json:"accessToken"`
	RefreshToken     string      `db:"refreshToken" json:"refreshToken"`
}

func GetCRESTTokens(characterID int64) ([]CRESTToken, error) {
	tokens := []CRESTToken{}
	if err := database.Select(&tokens, `
		SELECT characterID, tokenCharacterID, characterName,  expiry, tokenType, lastCode, lastStatus
		FROM crestTokens
		WHERE characterID = ?;`, characterID); err != nil {

		return nil, err
	}
	return tokens, nil
}

func AddCRESTToken(characterID int64, tokenCharacterID int64, characterName string, tok *eveapi.CRESTToken) error {
	if _, err := database.Exec(`
		INSERT INTO crestTokens	(characterID, tokenCharacterID, accessToken, refreshToken, expiry, tokenType, characterName)
			VALUES		(?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE 
				accessToken 	= VALUES(accessToken),
				refreshToken 	= VALUES(refreshToken),
				expiry 			= VALUES(expiry),
				tokenType 		= VALUES(tokenType)`,
		characterID, tokenCharacterID, tok.AccessToken, tok.RefreshToken, tok.Expiry, tok.TokenType, characterName); err != nil {

		return err
	}
	return nil
}

func DeleteCRESTToken(characterID int64, tokenCharacterID int64) error {
	if _, err := database.Exec(`DELETE FROM crestTokens WHERE characterID = ? AND tokenCharacterID = ? LIMIT 1`,
		characterID, tokenCharacterID); err != nil {

		return err
	}
	return nil
}

func UpdateCharacter(characterID int64, name string, bloodlineID int64, ancestryID int64, corporationID int64, allianceID int64,
	race string, securityStatus float64, cacheUntil time.Time) error {
	cacheUntil = time.Now().UTC().Add(time.Hour * 24 * 5)
	if _, err := database.Exec(`
		INSERT INTO characters (characterID,name,bloodlineID,ancestryID,corporationID,allianceID,race,securityStatus,updated,cacheUntil)
			VALUES(?,?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
			corporationID=VALUES(corporationID), allianceID=VALUES(allianceID), securityStatus=VALUES(securityStatus), updated = UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, characterID, name, bloodlineID, ancestryID, corporationID, allianceID, race, securityStatus, cacheUntil); err != nil {
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
		
		FROM characters C
		INNER JOIN corporations Co ON Co.corporationID = C.corporationID
		LEFT OUTER JOIN alliances Al ON Al.allianceID = C.allianceID
		WHERE characterID = ?
		LIMIT 1`, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}
