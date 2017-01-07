package models

import "time"

// Update an alliances information.
func UpdateAlliance(allianceID int64, name string, memberCount int64, shortName string, executorCorp int64,
	startDate time.Time, deleted bool, description string, creatorCorp int64, creatorCharacter int64,
	cacheUntil time.Time) error {

	cacheUntil = time.Now().UTC().Add(time.Hour * 24 * 1)
	if _, err := database.Exec(`
		INSERT INTO evedata.alliances 
			(
				allianceID,
				name,
				shortName,
				executorCorpID,
				startDate,
				corporationsCount,
				deleted,
				description,
				creatorCorpID,
				creatorCharacter,
				updated,
				cacheUntil
			)
			VALUES(?,?,?,?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
				executorCorpID = VALUES(executorCorpID),
				corporationsCount = VALUES(corporationsCount), 
				description = VALUES(description), 
				deleted = VALUES(deleted), 
				updated = UTC_TIMESTAMP(), 
				cacheUntil=VALUES(cacheUntil)
	`, allianceID, name, shortName, executorCorp, startDate, memberCount, deleted, description,
		creatorCorp, creatorCharacter, cacheUntil); err != nil {
		return err
	}
	return nil
}

type Alliance struct {
	AllianceID              int64     `db:"allianceID" json:"allianceID"`
	AllianceName            string    `db:"allianceName" json:"allianceName"`
	AllianceTicker          string    `db:"allianceTicker" json:"allianceTicker"`
	Description             string    `db:"description" json:"description"`
	CorporationsCount       int64     `db:"corporationsCount" json:"corporationsCount"`
	StartDate               time.Time `db:"startDate" json:"startDate"`
	ExecutorCorporationID   int64     `db:"executorCorporationID" json:"executorCorporationID"`
	ExecutorCorporationName string    `db:"executorCorporationName" json:"executorCorporationName"`
}

// Obtain alliance information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetAlliance(id int64) (*Alliance, error) {
	ref := Alliance{}
	if err := database.QueryRowx(`
		SELECT 
			A.allianceID,
		    A.name AS allianceName, 
		    A.shortName AS allianceTicker,
		    A.description,
		    A.corporationsCount,
		    A.startDate,
		    
		    EXEC.name AS executorCorporationName,
		    EXEC.corporationID AS executorCorporationID
		    
		FROM evedata.alliances A
		INNER JOIN corporations EXEC ON A.executorCorpID = EXEC.corporationID
		WHERE A.allianceID = ?
		LIMIT 1`, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

type AllianceMember struct {
	CorporationID   int64  `db:"corporationID" json:"corporationID"`
	CorporationName string `db:"corporationName" json:"corporationName"`
	MemberCount     int64  `db:"memberCount" json:"memberCount"`
}

// Obtain a list of corporations within an alliance by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetAllianceMembers(id int64) ([]AllianceMember, error) {
	ref := []AllianceMember{}
	if err := database.Select(&ref, `
		SELECT 
			M.corporationID, 
		    name AS corporationName,
		    M.memberCount
		FROM corporations M
		WHERE allianceID = ?;
		`, id); err != nil {
		return nil, err
	}
	return ref, nil
}
