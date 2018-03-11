package models

import (
	"errors"

	"github.com/antihax/evedata/services/conservator"
)

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetShares(characterID int32) ([]conservator.Share, error) {
	shares := []conservator.Share{}
	if err := database.Select(&shares, `
		SELECT S.characterID, S.tokenCharacterID, characterName AS tokenCharacterName, entityID, types, IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS entityType
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
func GetBotServices(characterID int32) ([]conservator.Service, error) {
	services := []conservator.Service{}
	if err := database.Select(&services, `
		SELECT  S.botServiceID, S.name, entityID, address,  type, services, options,
		IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS entityType
				FROM evedata.botServices S
				LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
				LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
				LEFT OUTER JOIN evedata.botDelegate D ON D.botServiceID = S.botServiceID
				LEFT OUTER JOIN evedata.crestTokens T ON FIND_IN_SET("Director", T.roles) AND 
					T.characterID = ? AND (A.executorCorpID = T.corporationID OR 
					(T.allianceID = 0 AND T.corporationID = S.entityID))
				WHERE D.characterID = ? OR T.characterID = ?
				GROUP BY botServiceID;`, characterID, characterID, characterID); err != nil {
		return nil, err
	}
	return services, nil
}

type BotServiceDetails struct {
	conservator.Service
	Channels []conservator.Channel
	Shares   []conservator.Share
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetBotServiceDetails(characterID, serverID int32) (BotServiceDetails, error) {

	// let this perform our authorization checks
	service := BotServiceDetails{}
	row, err := database.Queryx(`
		SELECT  S.botServiceID, S.name, entityID, address,  type, services, options,
		IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS entityType
				FROM evedata.botServices S
				LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
				LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
				LEFT OUTER JOIN evedata.botDelegate D ON D.botServiceID = S.botServiceID
				LEFT OUTER JOIN evedata.crestTokens T ON FIND_IN_SET("Director", T.roles) AND 
					T.characterID = ? AND (A.executorCorpID = T.corporationID OR 
					(T.allianceID = 0 AND T.corporationID = S.entityID))
				WHERE (D.characterID = ? OR T.characterID = ?) AND S.botServiceID = ?
				GROUP BY botServiceID LIMIT 1;`, characterID, characterID, characterID, serverID)
	if err != nil {
		return service, err
	}
	defer row.Close()

	if !row.Next() {
		return service, errors.New("Bot service unavailable")
	}
	err = row.StructScan(&service)
	if err != nil {
		return service, err
	}

	// not authorized
	if service.BotServiceID == 0 {
		return service, errors.New("character is not authorized")
	}

	err = database.Select(&service.Channels, `
		SELECT botServiceID, channelID, channelName, services, options
		FROM evedata.botChannels 
		WHERE botServiceID = ?;`, service.BotServiceID)
	if err != nil {
		return service, err
	}

	err = database.Select(&service.Shares, `
		SELECT characterName AS tokenCharacterName, C.tokenCharacterID, E.corporationID AS entityID, E.name AS entityName, "corporation" AS entityType, B.botServiceID, types, ignored
		FROM evedata.sharing S
		INNER JOIN evedata.botServices B ON B.entityID = S.entityID
        INNER JOIN evedata.crestTokens C ON C.tokenCharacterID = S.tokenCharacterID
        INNER JOIN evedata.corporations E ON C.corporationID = E.corporationID
		WHERE B.botServiceID = ?;`, service.BotServiceID)
	if err != nil {
		return service, err
	}

	return service, nil
}

func AddBotServiceChannel(botServiceID int32, channelID, channelName string) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.botChannels	(botServiceID, channelID, channelName, services, options)
			VALUES(?,?,?,'','')
			ON DUPLICATE KEY UPDATE channelID = channelID`,
		botServiceID, channelID, channelName); err != nil {
		return err
	}

	return nil
}

func DeleteBotServiceChannel(botServiceID int32, channelID string) error {
	if _, err := database.Exec(`
		DELETE  FROM evedata.botChannels
		WHERE botServiceID = ? AND channelID = ?`, botServiceID, channelID); err != nil {
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
		INSERT INTO evedata.botServices	(entityID, address, type, options)
			VALUES(?,?,'discord','')
			ON DUPLICATE KEY UPDATE entityID = entityID`,
		entityID, serverID); err != nil {
		return err
	}

	return nil
}

func DeleteService(characterID, botServiceID int32) error {
	if _, err := database.Exec(`
		DELETE S FROM evedata.botServices S
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
		LEFT OUTER JOIN evedata.crestTokens T ON FIND_IN_SET("Director", T.roles) AND 
		T.characterID = ? AND (A.executorCorpID = T.corporationID OR 
		(T.allianceID = 0 AND T.corporationID = S.entityID))
		WHERE botServiceID = ? AND T.characterID = ?
		`,
		characterID, botServiceID, characterID); err != nil {
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
