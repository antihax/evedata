package models

type KnownAlts struct {
	CharacterID   int64  `db:"characterID" json:"id"`
	CharacterName string `db:"characterName" json:"name"`
	Frequency     int    `db:"frequency" json:"frequency"`
	Type          string `db:"type" json:"type"`
	Source        uint8  `db:"source" json:"source"`
}

// Obtain Character Associates by ID.

func GetCharacterKnownAssociates(id int64) ([]KnownAlts, error) {
	ref := []KnownAlts{}
	if err := database.Select(&ref, `
		SELECT 	associateID AS characterID,
				frequency,
				C.name AS characterName,
				IFNULL(source, 0) AS source
		FROM evedata.characterAssociations A
		INNER JOIN evedata.characters C ON A.associateID = C.characterID
		INNER JOIN evedata.characters M ON A.characterID = M.characterID
		WHERE A.characterID = ?
		AND (M.allianceID != C.allianceID OR C.allianceID = 0) AND M.corporationID != C.corporationID
		`, id); err != nil {
		return nil, err
	}

	for i := range ref {
		ref[i].Type = "character"
	}

	return ref, nil
}

func GetCorporationKnownAssociates(id int64) ([]KnownAlts, error) {
	ref := []KnownAlts{}
	if err := database.Select(&ref, `
		SELECT	associateID AS characterID,
				SUM(frequency) AS frequency,
				C.name AS characterName,
				IFNULL(source, 0) AS source
		FROM evedata.characterAssociations A
        INNER JOIN evedata.characters C ON A.associateID = C.characterID
        INNER JOIN evedata.characters M ON A.characterID = M.characterID
        WHERE A.characterID IN (SELECT characterID FROM evedata.characters WHERE corporationID = ?)
        AND (M.allianceID != C.allianceID OR C.allianceID = 0) AND M.corporationID != C.corporationID
        GROUP BY associateID
		`, id); err != nil {
		return nil, err
	}

	for i := range ref {
		ref[i].Type = "character"
	}

	return ref, nil
}

func GetAllianceKnownAssociates(id int64) ([]KnownAlts, error) {
	ref := []KnownAlts{}
	if err := database.Select(&ref, `
		SELECT	associateID AS characterID,
				SUM(frequency) AS frequency,
				C.name AS characterName,
				IFNULL(source, 0) AS source
		FROM evedata.characterAssociations A
        INNER JOIN evedata.characters C ON A.associateID = C.characterID
        INNER JOIN evedata.characters M ON A.characterID = M.characterID
        WHERE A.characterID IN (SELECT characterID FROM evedata.characters WHERE allianceID = ?)
        AND (M.allianceID != C.allianceID OR C.allianceID = 0) AND M.corporationID != C.corporationID
        GROUP BY associateID
		`, id); err != nil {
		return nil, err
	}

	for i := range ref {
		ref[i].Type = "character"
	}

	return ref, nil
}
