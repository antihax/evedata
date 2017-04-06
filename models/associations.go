package models

type KnownAlts struct {
	CharacterID   int64  `db:"characterID" json:"id"`
	CharacterName string `db:"characterName" json:"name"`
	Frequency     int    `db:"frequency" json:"frequency"`
	Type          string `db:"type" json:"type"`
}

// Obtain Character information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCharacterKnownAssociates(id int64) ([]KnownAlts, error) {
	ref := []KnownAlts{}
	if err := database.Select(&ref, `
		SELECT 
			C.characterID,
			C.name AS characterName,
			frequency
		    
		FROM evedata.characterAssociations A
		INNER JOIN evedata.characters C ON C.characterID = A.associateID
		WHERE A.characterID = ?
		`, id); err != nil {
		return nil, err
	}

	for i := range ref {
		ref[i].Type = "character"
	}

	return ref, nil
}

// Obtain Character information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCharacterKnownKillmailAssociates(id int64) ([]KnownAlts, error) {
	ref := []KnownAlts{}
	if err := database.Select(&ref, `
		SELECT 	associateID AS characterID,
				frequency,
				C.name AS characterName
		FROM evedata.characterKillmailAssociations A
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
