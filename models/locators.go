package models

type LocatorShares struct {
	CharacterID int64  `db:"characterID" json:"characterID,omitempty"`
	ID          int64  `db:"entityID" json:"id,omitempty"`
	EntityName  string `db:"entityName" json:"entityName,omitempty"`
	EntityType  string `db:"type" json:"type,omitempty"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetLocatorShares(characterID int64) ([]LocatorShares, error) {
	locatorShares := []LocatorShares{}
	if err := database.Select(&locatorShares, `
		SELECT characterID, entityID, IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS type
		FROM evedata.locatorShareWith S
		LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
		WHERE characterID = ?;`, characterID); err != nil {
		return nil, err
	}
	return locatorShares, nil
}

func AddLocatorShare(characterID int64, entityID int64) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.locatorShareWith	(characterID, entityID)
			VALUES(?,?)
			ON DUPLICATE KEY UPDATE entityID = entityID"`,
		characterID, entityID); err != nil {
		return err
	}
	return nil
}

func DeleteLocatorShare(characterID int64, entityID int64) error {
	if _, err := database.Exec(`DELETE FROM evedata.locatorShareWith WHERE characterID = ? AND entityID = ? LIMIT 1`,
		characterID, entityID); err != nil {
		return err
	}
	return nil
}
