package models

import (
	"time"

	"github.com/antihax/evedata/eveapi"
	"github.com/antihax/evedata/null"
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

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCRESTTokens(characterID int64) ([]CRESTToken, error) {
	tokens := []CRESTToken{}
	if err := database.Select(&tokens, `
		SELECT characterID, tokenCharacterID, characterName,  expiry, tokenType, lastCode, lastStatus
		FROM evedata.crestTokens
		WHERE characterID = ?;`, characterID); err != nil {

		return nil, err
	}
	return tokens, nil
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

func AddCRESTToken(characterID int64, tokenCharacterID int64, characterName string, tok *eveapi.CRESTToken, scopes string) error {
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

func UpdateCharacter(characterID int64, name string, bloodlineID int64, ancestryID int64, corporationID int64, allianceID int64,
	race string, securityStatus float64, cacheUntil time.Time) error {
	cacheUntil = time.Now().UTC().Add(time.Hour * 24 * 5)
	if _, err := database.Exec(`
		INSERT INTO evedata.characters (characterID,name,bloodlineID,ancestryID,corporationID,allianceID,race,securityStatus,updated,cacheUntil)
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
