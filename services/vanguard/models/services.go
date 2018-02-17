package models

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

type Service struct {
	BotServiceID int32  `db:"botServiceID" json:"botServiceID"`
	Name         string `db:"name" json:"name"`
	EntityID     int32  `db:"entityID" json:"entityID"`
	EntityName   string `db:"entityName" json:"entityName"`
	EntityType   string `db:"entityType" json:"entityType"`
	Address      string `db:"address" json:"address" `
	Type         string `db:"type" json:"type"`
	Services     string `db:"services" json:"services"`
	OptionsJSON  string `db:"options" json:"options"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetBotServices(characterID int32) ([]Service, error) {
	services := []Service{}
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

func AddService(characterID, tokenCharacterID, entityID int32, types string) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.sharing	(characterID, tokenCharacterID, entityID, types)
			VALUES(?,?,?,?)
			ON DUPLICATE KEY UPDATE entityID = entityID, types = VALUES(types)`,
		characterID, tokenCharacterID, entityID, types); err != nil {
		return err
	}
	return nil
}

func DeleteService(characterID, botServiceID int32) error {
	if _, err := database.Exec(`
		DELETE FROM evedata.botServices S
		LEFT OUTER JOIN evedata.crestTokens T ON FIND_IN_SET("Director", T.roles) AND 
		T.characterID = ? AND (A.executorCorpID = T.corporationID OR 
		(T.allianceID = 0 AND T.corporationID = S.entityID))
		WHERE botServiceID = ? AND T.characterID = ?
		LIMIT 1`,
		characterID, botServiceID, characterID); err != nil {
		return err
	}
	return nil
}
