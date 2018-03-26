package models

import (
	"errors"
	"time"

	"github.com/antihax/evedata/services/conservator"
)

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetShares(characterID int32) ([]conservator.Share, error) {
	shares := []conservator.Share{}
	if err := database.Select(&shares, `
		SELECT DISTINCT S.characterID, S.tokenCharacterID, characterName AS tokenCharacterName, entityID, types, IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS entityType
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

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetIntegrations(characterID int32) ([]conservator.Service, error) {
	services := []conservator.Service{}
	if err := database.Select(&services, `
		SELECT  S.integrationID, S.name, entityID, address,  type, services, options,
		IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS entityType
				FROM evedata.integrations S
				LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
				LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
				LEFT OUTER JOIN evedata.integrationDelegate D ON D.integrationID = S.integrationID
				LEFT OUTER JOIN evedata.crestTokens T ON FIND_IN_SET("Director", T.roles) AND 
					T.characterID = ? AND (A.executorCorpID = T.corporationID OR 
					(T.allianceID = 0 AND T.corporationID = S.entityID))
				WHERE D.characterID = ? OR T.characterID = ?
				GROUP BY integrationID;`, characterID, characterID, characterID); err != nil {
		return nil, err
	}
	return services, nil
}

type IntegrationDetails struct {
	conservator.Service
	Channels []conservator.Channel
	Shares   []conservator.Share
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetIntegrationDetails(characterID, serverID int32) (IntegrationDetails, error) {
	// let this perform our authorization checks
	service := IntegrationDetails{}
	row, err := database.Queryx(`
		SELECT  S.integrationID, S.name, entityID, address,  type, services, options,
		IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS entityType
				FROM evedata.integrations S
				LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
				LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
				LEFT OUTER JOIN evedata.integrationDelegate D ON D.integrationID = S.integrationID
				LEFT OUTER JOIN evedata.crestTokens T ON FIND_IN_SET("Director", T.roles) AND 
					T.characterID = ? AND (A.executorCorpID = T.corporationID OR 
					(T.allianceID = 0 AND T.corporationID = S.entityID))
				WHERE (D.characterID = ? OR T.characterID = ?) AND S.integrationID = ?
				GROUP BY integrationID LIMIT 1;`, characterID, characterID, characterID, serverID)
	if err != nil {
		return service, err
	}
	defer row.Close()

	if !row.Next() {
		return service, errors.New("Integration service unavailable")
	}
	err = row.StructScan(&service)
	if err != nil {
		return service, err
	}

	// not authorized
	if service.IntegrationID == 0 {
		return service, errors.New("character is not authorized")
	}

	err = database.Select(&service.Channels, `
		SELECT integrationID, channelID, channelName, services, options
		FROM evedata.integrationChannels 
		WHERE integrationID = ?;`, service.IntegrationID)
	if err != nil {
		return service, err
	}

	err = database.Select(&service.Shares, `
		SELECT DISTINCT characterName AS tokenCharacterName, C.tokenCharacterID, E.corporationID AS entityID, E.name AS entityName, "corporation" AS entityType, B.integrationID, types, ignored
		FROM evedata.sharing S
		INNER JOIN evedata.integrations B ON B.entityID = S.entityID
        INNER JOIN evedata.crestTokens C ON C.tokenCharacterID = S.tokenCharacterID
        INNER JOIN evedata.corporations E ON C.corporationID = E.corporationID
		WHERE B.integrationID = ?;`, service.IntegrationID)
	if err != nil {
		return service, err
	}

	return service, nil
}

func AddIntegrationChannel(integrationID int32, channelID, channelName string) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.integrationChannels	(integrationID, channelID, channelName, services, options)
			VALUES(?,?,?,'','')
			ON DUPLICATE KEY UPDATE channelID = channelID`,
		integrationID, channelID, channelName); err != nil {
		return err
	}

	return nil
}

func DeleteIntegrationChannel(integrationID int32, channelID string) error {
	if _, err := database.Exec(`
		DELETE FROM evedata.integrationChannels
		WHERE integrationID = ? AND channelID = ?`, integrationID, channelID); err != nil {
		return err
	}
	return nil
}

func AddDiscordService(characterID, entityID int32, serverID string) error {
	// verify this user is able to create a discord service for this entity
	entities, err := GetEntitiesWithRole(characterID, "Director")
	if err != nil {
		return err
	}

	if !entityInSlice(entityID, entities) {
		return errors.New("character is unauthorized to create this discord entry")
	}

	if _, err := database.Exec(`
		INSERT INTO evedata.integrations	(entityID, address, type, options)
			VALUES(?,?,'discord','')
			ON DUPLICATE KEY UPDATE entityID = entityID`,
		entityID, serverID); err != nil {
		return err
	}

	return nil
}

func DeleteService(characterID, integrationID int32) error {
	if _, err := database.Exec(`
		DELETE S FROM evedata.integrations S
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
		LEFT OUTER JOIN evedata.crestTokens T ON FIND_IN_SET("Director", T.roles) AND 
		T.characterID = ? AND (A.executorCorpID = T.corporationID OR 
		(T.allianceID = 0 AND T.corporationID = S.entityID))
		WHERE integrationID = ? AND T.characterID = ?
		`,
		characterID, integrationID, characterID); err != nil {
		return err
	}
	return nil
}

func UpdateChannel(integrationID int32, channelID, options, services string) error {
	if _, err := database.Exec(`
		UPDATE evedata.integrationChannels SET options = ?, services = ? WHERE integrationID = ? AND channelID = ? LIMIT 1`,
		options, services, integrationID, channelID); err != nil {
		return err
	}
	return nil
}

func UpdateService(integrationID int32, options, services string) error {
	if _, err := database.Exec(`
		UPDATE evedata.integrations SET options = ?, services = ? WHERE integrationID = ? LIMIT 1`,
		options, services, integrationID); err != nil {
		return err
	}
	return nil
}

func entityInSlice(a int32, list []Entity) bool {
	for _, b := range list {
		if b.EntityID == a {
			return true
		}
	}
	return false
}

type AvailableIntegrations struct {
	IntegrationID     int32     `db:"integrationID" json:"integrationID"`
	Address           string    `db:"address" json:"address"`
	Reason            string    `db:"reason" json:"reason"`
	Name              string    `db:"name" json:"name"`
	CharacterName     string    `db:"characterName" json:"characterName" `
	CharacterID       int32     `db:"characterID" json:"characterID"`
	TokenCharacterID  int32     `db:"tokenCharacterID" json:"tokenCharacterID"`
	IntegrationUserID string    `db:"integrationUserID" json:"integrationUserID"`
	Type              string    `db:"type" json:"type"`
	EntityID          int32     `db:"entityID" json:"entityID"`
	EntityName        string    `db:"entityName" json:"entityName"`
	EntityType        string    `db:"entityType" json:"entityType"`
	Expiry            time.Time `db:"expiry" json:"expiry,omitempty"`
	AccessToken       string    `db:"accessToken" json:"accessToken,omitempty"`
	RefreshToken      string    `db:"refreshToken" json:"refreshToken,omitempty"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetAvailableIntegrations(characterID int32) ([]AvailableIntegrations, error) {
	integrations := []AvailableIntegrations{}
	if err := database.Select(&integrations, `
		SELECT integrationID, address, reason, S.name, characterName, T.characterID, tokenCharacterID, 
			integrationUserID, type, entityID, IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS entityType
		FROM
		(
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "member" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.crestTokens C ON C.authCharacter = 1 AND
			(C.corporationID = B.entityID 					   
			OR C.allianceID = B.entityID)
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%member%"
		UNION
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "militia" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.crestTokens C ON C.authCharacter = 1 AND
			B.factionID > 0 AND B.factionID = C.factionID
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%militia%"
		UNION
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "alliedMilitia" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.crestTokens C ON C.authCharacter = 1 AND
			B.factionID > 0 AND B.factionID = evedata.alliedMilita(C.factionID)
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%alliedMilitia%"
		UNION
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "+5" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.entityContacts E ON E.entityID = B.entityID AND E.standing = 5.0
		INNER JOIN evedata.crestTokens C ON  C.authCharacter = 1 AND
			(E.contactID = C.tokenCharacterID 
			OR E.contactID = C.corporationID
			OR E.contactID = C.allianceID)
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%plusFive%"
		UNION
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "+10" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.entityContacts E ON E.entityID = B.entityID AND E.standing = 10.0
		INNER JOIN evedata.crestTokens C ON C.authCharacter = 1 AND
			(E.contactID = C.tokenCharacterID 
			OR E.contactID = C.corporationID
			OR E.contactID = C.allianceID)
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%plusTen%"
		) S 
		INNER JOIN evedata.integrationTokens T ON S.characterID = T.characterID
		LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
		WHERE T.characterID = ?
		GROUP BY address `, characterID); err != nil {
		return nil, err
	}
	return integrations, nil
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetIntegrationsForCharacter(characterID, integrationID int32) (*AvailableIntegrations, error) {
	integration := []AvailableIntegrations{}
	if err := database.Select(&integration, `	
		SELECT integrationID, address, accessToken, refreshToken, expiry, S.name, characterName, T.characterID, tokenCharacterID, 
		integrationUserID, type, entityID, IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS entityType
		FROM
		(
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "member" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.crestTokens C ON C.authCharacter = 1 AND
		(C.corporationID = B.entityID 					   
		OR C.allianceID = B.entityID)
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%member%"
		UNION
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "militia" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.crestTokens C ON C.authCharacter = 1 AND
		B.factionID > 0 AND B.factionID = C.factionID
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%militia%"
		UNION
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "alliedMilitia" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.crestTokens C ON C.authCharacter = 1 AND
		B.factionID > 0 AND B.factionID = evedata.alliedMilita(C.factionID)
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%alliedMilitia%"
		UNION
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "+5" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.entityContacts E ON E.entityID = B.entityID AND E.standing = 5.0
		INNER JOIN evedata.crestTokens C ON  C.authCharacter = 1 AND
		(E.contactID = C.tokenCharacterID 
		OR E.contactID = C.corporationID
		OR E.contactID = C.allianceID)
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%plusFive%"
		UNION
		SELECT integrationID, address, B.entityID, name, characterName, C.characterID, tokenCharacterID, "+10" AS reason
		FROM evedata.integrations B
		INNER JOIN evedata.entityContacts E ON E.entityID = B.entityID AND E.standing = 10.0
		INNER JOIN evedata.crestTokens C ON C.authCharacter = 1 AND
		(E.contactID = C.tokenCharacterID 
		OR E.contactID = C.corporationID
		OR E.contactID = C.allianceID)
		WHERE FIND_IN_SET(B.services, "auth") AND options LIKE "%plusTen%"
		) S 
		INNER JOIN evedata.integrationTokens T ON S.characterID = T.characterID
		LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
		WHERE T.characterID = ? AND integrationID = ?
		GROUP BY address `, characterID, integrationID); err != nil {
		return nil, err
	}
	return &integration[0], nil
}
