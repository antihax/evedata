package models

import "time"

// Update an alliances information.
func UpdateAlliance(allianceID int32, name string, memberCount int, shortName string, executorCorp int32,
	startDate time.Time, cacheUntil time.Time) error {

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
				updated,
				cacheUntil
			)
			VALUES(?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
				executorCorpID = VALUES(executorCorpID),
				corporationsCount = VALUES(corporationsCount), 
				updated = UTC_TIMESTAMP(), 
				cacheUntil=VALUES(cacheUntil)
	`, allianceID, name, shortName, executorCorp, startDate, memberCount, cacheUntil); err != nil {
		return err
	}
	return nil
}

type Alliance struct {
	AllianceID              int64     `db:"allianceID" json:"allianceID"`
	AllianceName            string    `db:"allianceName" json:"allianceName"`
	AllianceTicker          string    `db:"allianceTicker" json:"allianceTicker"`
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
		    A.corporationsCount,
		    A.startDate,
		    
		    EXEC.name AS executorCorporationName,
		    EXEC.corporationID AS executorCorporationID
		    
		FROM evedata.alliances A
		INNER JOIN evedata.corporations EXEC ON A.executorCorpID = EXEC.corporationID
		WHERE A.allianceID = ?
		LIMIT 1`, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

type AllianceMember struct {
	ID              int64  `db:"corporationID" json:"id"`
	CorporationName string `db:"corporationName" json:"name"`
	MemberCount     int64  `db:"memberCount" json:"memberCount"`
	Type            string `db:"type" json:"type"`
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
		FROM evedata.corporations M
		WHERE allianceID = ?;
		`, id); err != nil {
		return nil, err
	}

	for i := range ref {
		ref[i].Type = "corporation"
	}

	return ref, nil
}
