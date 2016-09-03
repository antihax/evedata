package models

import (
	"evedata/eveapi"
	"evedata/null"
	"time"
)

type ApiKey struct {
	KeyID       int         `db:"keyID" json:"keyID"`
	CharacterID int         `db:"characterID" json:"characterID"`
	NextCheck   time.Time   `db:"nextCheck" json:"nextCheck"`
	LastCode    int         `db:"lastCode" json:"lastCode"`
	LastError   null.String `db:"lastError" json:"lastError"`
	AccessMask  int         `db:"accessMask" json:"accessMask"`
	Type        null.String `db:"type" json:"type"`
}

type CRESTToken struct {
	Expiry           time.Time   `db:"expiry" json:"expiry"`
	CharacterID      int         `db:"characterID" json:"characterID"`
	TokenType        string      `db:"tokenType" json:"tokenType"`
	TokenCharacterID int         `db:"tokenCharacterID" json:"tokenCharacterID"`
	CharacterName    string      `db:"characterName" json:"characterName"`
	LastCode         int         `db:"lastCode" json:"lastCode"`
	LastStatus       null.String `db:"lastStatus" json:"lastStatus"`
	AccessToken      string      `db:"accessToken" json:"accessToken"`
	RefreshToken     string      `db:"refreshToken" json:"refreshToken"`
}

func (c *ApiKey) UpdateChecked(nextCheck time.Time) error {
	if _, err := database.Exec(`
		UPDATE apiKeys SET nextCheck = ? 
		WHERE characterID = ? AND keyID = ? 
		LIMIT 1;`, nextCheck, c.CharacterID, c.KeyID); err != nil {

		return err
	}
	return nil
}

func GetAPIKeys(characterID int64) ([]ApiKey, error) {
	keys := []ApiKey{}
	if err := database.Select(&keys, `
		SELECT keyID, characterID, nextCheck, accessMask, lastCode, lastError, type 
		FROM apiKeys 
		WHERE characterID = ?;`, characterID); err != nil {

		return nil, err
	}
	return keys, nil
}

func AddApiKey(characterID int64, keyID int, vCode string) error {
	if _, err := database.Exec(`INSERT INTO apiKeys (characterID, keyID, vCode)VALUES(?,?,?)`,
		characterID, keyID, vCode); err != nil {

		return err
	}
	return nil
}

func DeleteApiKey(characterID int64, keyID int) error {
	if _, err := database.Exec(`DELETE FROM apiKeys WHERE characterID = ? AND keyID = ? LIMIT 1`,
		characterID, keyID); err != nil {

		return err
	}
	return nil
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

func DeleteCRESTToken(characterID int64, tokenCharacterID int) error {
	if _, err := database.Exec(`DELETE FROM crestTokens WHERE characterID = ? AND tokenCharacterID = ? LIMIT 1`,
		characterID, tokenCharacterID); err != nil {

		return err
	}
	return nil
}

func UpdateCharacter(characterID int64, name string, bloodlineID int64, ancestryID int64, corporationID int64, allianceID int64,
	race string, securityStatus float64, cacheUntil time.Time) error {
	cacheUntil = time.Now().UTC().Add(time.Hour * 24)
	if _, err := database.Exec(`
		INSERT INTO eve.character (characterID,name,bloodlineID,ancestryID,corporationID,allianceID,race,securityStatus,updated,cacheUntil)
			VALUES(?,?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
			corporationID=VALUES(corporationID), allianceID=VALUES(allianceID), securityStatus=VALUES(securityStatus), updated = UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, characterID, name, bloodlineID, ancestryID, corporationID, allianceID, race, securityStatus, cacheUntil); err != nil {
		return err
	}
	return nil
}
