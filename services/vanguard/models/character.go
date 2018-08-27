package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/conservator"
	"github.com/guregu/null"
	"golang.org/x/oauth2"
)

// Obtain an authenticated client from a stored access/refresh token.
func GetCRESTToken(characterID int32, ownerHash string, tokenCharacterID int32) (*CRESTToken, error) {
	tok := &CRESTToken{}
	if err := database.QueryRowx(
		`SELECT expiry, tokenType, accessToken, refreshToken, tokenCharacterID, characterID, characterName, corporationID, allianceID
			FROM evedata.crestTokens
			WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = "") AND tokenCharacterID = ?
			LIMIT 1`,
		characterID, ownerHash, tokenCharacterID).StructScan(tok); err != nil {

		return nil, err
	}

	return tok, nil
}

type CRESTToken struct {
	Expiry           time.Time           `db:"expiry" json:"expiry,omitempty"`
	CharacterID      int32               `db:"characterID" json:"characterID,omitempty"`
	CorporationID    int32               `db:"corporationID" json:"corporationID,omitempty"`
	AllianceID       int32               `db:"allianceID" json:"allianceID,omitempty"`
	TokenType        string              `db:"tokenType" json:"tokenType,omitempty"`
	TokenCharacterID int32               `db:"tokenCharacterID" json:"tokenCharacterID,omitempty"`
	CharacterName    string              `db:"characterName" json:"characterName,omitempty"`
	CorporationName  null.String         `db:"corporationName" json:"corporationName,omitempty"`
	AllianceName     null.String         `db:"allianceName" json:"allianceName,omitempty"`
	LastCode         int64               `db:"lastCode" json:"lastCode,omitempty"`
	LastStatus       null.String         `db:"lastStatus" json:"lastStatus,omitempty"`
	AccessToken      string              `db:"accessToken" json:"accessToken,omitempty"`
	RefreshToken     string              `db:"refreshToken" json:"refreshToken,omitempty"`
	Scopes           string              `db:"scopes" json:"scopes"`
	AuthCharacter    int                 `db:"authCharacter" json:"authCharacter"`
	SharingInt       string              `db:"sharingint" json:"_,omitempty"`
	Sharing          []conservator.Share `json:"sharing"`
	MailPassword     int                 `db:"mailPassword" json:"mailPassword"`
}

type IntegrationToken struct {
	Type                string      `db:"type" json:"type,omitempty"`
	Expiry              time.Time   `db:"expiry" json:"expiry,omitempty"`
	CharacterID         int32       `db:"characterID" json:"characterID,omitempty"`
	IntegrationUserID   string      `db:"integrationUserID" json:"integrationUserID,omitempty"`
	IntegrationUserName string      `db:"integrationUserName" json:"integrationUserName,omitempty"`
	TokenType           string      `db:"tokenType" json:"tokenType,omitempty"`
	LastCode            int64       `db:"lastCode" json:"lastCode,omitempty"`
	LastStatus          null.String `db:"lastStatus" json:"lastStatus,omitempty"`
	AccessToken         string      `db:"accessToken" json:"accessToken,omitempty"`
	RefreshToken        string      `db:"refreshToken" json:"refreshToken,omitempty"`
	Scopes              string      `db:"scopes" json:"scopes"`
}

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

func GetCRESTTokens(characterID int32, ownerHash string) ([]CRESTToken, error) {
	tokens := []CRESTToken{}
	if err := database.Select(&tokens, `
		SELECT T.characterID, T.tokenCharacterID, characterName, IF(mailPassword != "", 1, 0) AS mailPassword,
		lastCode, lastStatus, scopes, authCharacter, C1.name AS corporationName, A1.name AS allianceName,
		T.corporationID, T.allianceID,
		IFNULL(
			CONCAT("[", GROUP_CONCAT(CONCAT(
				'{"id": ', entityID, 
				', "types": "', types, '"',
				', "entityName": "', IFNULL(A.name, C.name), '"',
				', "type": "', IF(A.name IS NULL, "corporation", "alliance"), '"',
				'}')), 
			"]")
		, "[]") AS sharingint
		FROM evedata.crestTokens T
		LEFT OUTER JOIN evedata.sharing S ON T.tokenCharacterID = S.tokenCharacterID AND T.characterID = S.characterID
		LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
		LEFT OUTER JOIN evedata.corporations C1 ON C1.corporationID = T.corporationID
		LEFT OUTER JOIN evedata.alliances A1 ON A1.allianceID = T.allianceID
		WHERE T.characterID = ? AND (T.characterOwnerHash = ? OR T.characterOwnerHash = "")
		GROUP BY characterID, tokenCharacterID;
		;`, characterID, ownerHash); err != nil {

		return nil, err
	}

	// Unmarshal our sharing data.
	for index := range tokens {
		share := []conservator.Share{}
		if err := json.Unmarshal([]byte(tokens[index].SharingInt), &share); err != nil {
			return nil, err
		}
		tokens[index].Sharing = share
		tokens[index].SharingInt = ""
	}
	return tokens, nil
}

// AddCRESTToken adds an SSO token to the database or updates it if one exists.
// resetting status and if errors were mailed to the user.
func AddCRESTToken(characterID int32, tokenCharacterID int32, characterName string, tok *oauth2.Token, scopes, ownerHash string, corporationID, allianceID, factionID int32) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.crestTokens	(characterID, tokenCharacterID, accessToken, refreshToken, expiry, 
				tokenType, characterName, scopes, lastStatus, characterOwnerHash, corporationID, allianceID, factionID)
			VALUES		(?,?,?,?,?,?,?,?,"Unused",?,?,?,?)
			ON DUPLICATE KEY UPDATE 
				accessToken 		= VALUES(accessToken),
				refreshToken 		= VALUES(refreshToken),
				expiry 				= VALUES(expiry),
				tokenType 			= VALUES(tokenType),
				characterOwnerHash	= VALUES(characterOwnerHash),
				scopes 				= VALUES(scopes),
				corporationID 		= VALUES(corporationID),
				allianceID	 		= VALUES(allianceID),
				factionID	 		= VALUES(factionID),
				lastStatus			= "Unused",
				mailedError 		= 0`,
		characterID, tokenCharacterID, tok.AccessToken, tok.RefreshToken, tok.Expiry, tok.TokenType, characterName, scopes, ownerHash, corporationID, allianceID, factionID); err != nil {
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

// AddIntegrationToken adds a oauth2 token to the database for integrations or updates it if one exists.
// resetting status and if errors were mailed to the user.
func AddIntegrationToken(tokenType string, characterID int32, userID string, userName string, tok *oauth2.Token, scopes string) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.integrationTokens	(type, characterID, integrationUserID, integrationUserName, accessToken, refreshToken, expiry, 
				tokenType, scopes, lastStatus)
			VALUES		(?,?,?,?,?,?,?,?,?,"Unused")
			ON DUPLICATE KEY UPDATE 
				accessToken 		= VALUES(accessToken),
				refreshToken 		= VALUES(refreshToken),
				expiry 				= VALUES(expiry),
				tokenType 			= VALUES(tokenType),
				scopes 				= VALUES(scopes),
				lastStatus			= "Unused",
				mailedError 		= 0`,
		tokenType, characterID, userID, userName, tok.AccessToken, tok.RefreshToken, tok.Expiry, tok.TokenType, scopes); err != nil {
		return err
	}
	return nil
}

func SetMailPassword(characterID, tokenCharacterID int32, ownerHash, password string) error {
	// BCrypt the password
	hash, err := sqlhelper.Hash(password)
	if err != nil {
		return err
	}

	if _, err := database.Exec(`UPDATE evedata.crestTokens
		SET mailPassword = ?
		WHERE characterID = ? AND tokenCharacterID = ? AND characterOwnerHash = ?;
		`, hash, characterID, tokenCharacterID, ownerHash); err != nil {
		return err
	}
	return nil
}

func GetIntegrationTokens(characterID int32) ([]IntegrationToken, error) {
	tokens := []IntegrationToken{}
	if err := database.Select(&tokens, `
		SELECT characterID,
			integrationUserID,
			type,
			integrationUserName,
			expiry,
			tokenType,
			lastCode,
			lastStatus,
			scopes
			FROM evedata.integrationTokens
			WHERE characterID = ?;
		`, characterID); err != nil {

		return nil, err
	}

	return tokens, nil
}

func DeleteIntegrationToken(tokenType string, characterID int32, integrationUserID string) error {
	if _, err := database.Exec(`DELETE FROM evedata.integrationTokens WHERE characterID = ? AND integrationUserID = ? AND type = ? LIMIT 1`,
		characterID, integrationUserID, tokenType); err != nil {
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
	Efficiency      float64     `db:"efficiency" json:"efficiency"`
	CapKills        int64       `db:"capKills" json:"capKills"`
	Kills           int64       `db:"kills" json:"kills"`
	Losses          int64       `db:"losses" json:"losses"`
}

// Obtain Character information by ID.

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
		    securityStatus,
			coalesce(efficiency, 0) AS efficiency,
			coalesce(capKills, 0) AS capKills,
			coalesce(kills, 0) AS kills,
			coalesce(losses, 0) AS losses
		FROM evedata.characters C
		LEFT OUTER JOIN evedata.corporations Co ON Co.corporationID = C.corporationID
		LEFT OUTER JOIN evedata.alliances Al ON Al.allianceID = C.allianceID
		LEFT OUTER JOIN evedata.entityKillStats S ON S.id = C.characterID
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

type Entity struct {
	EntityID   int32  `db:"entityID" json:"entityID"`
	EntityName string `db:"entityName" json:"entityName"`
	EntityType string `db:"entityType" json:"entityType"`
}

// GetEntitiesWithRole determine which corporation/alliance roles are available

func GetEntitiesWithRole(characterID int32, role string) ([]Entity, error) {
	ref := []Entity{}
	if err := database.Select(&ref, `
		SELECT DISTINCT C.corporationID AS entityID, name AS entityName, "corporation" AS entityType
		FROM evedata.crestTokens T
		INNER JOIN evedata.corporations C ON C.corporationID = T.corporationID
		WHERE FIND_IN_SET(?, T.roles) AND T.characterID = ?
        UNION
   		SELECT DISTINCT A.allianceID AS entityID, name AS entityName, "alliance" AS entityType
		FROM evedata.crestTokens T
		INNER JOIN evedata.alliances A ON A.allianceID = T.allianceID AND T.corporationID = A.executorCorpID
		WHERE FIND_IN_SET(?, T.roles) AND T.characterID = ?
		`, role, characterID, role, characterID); err != nil {
		return nil, err
	}
	return ref, nil
}
