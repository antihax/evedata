package models

import (
	"time"

	"github.com/guregu/null"
)

// UpdateAlliance Update an alliances information.
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

// Alliance resultset
type Alliance struct {
	AllianceID              int64     `db:"allianceID" json:"allianceID"`
	AllianceName            string    `db:"allianceName" json:"allianceName"`
	AllianceTicker          string    `db:"allianceTicker" json:"allianceTicker"`
	CorporationsCount       int64     `db:"corporationsCount" json:"corporationsCount"`
	StartDate               time.Time `db:"startDate" json:"startDate"`
	ExecutorCorporationID   int64     `db:"executorCorporationID" json:"executorCorporationID"`
	ExecutorCorporationName string    `db:"executorCorporationName" json:"executorCorporationName"`
	Efficiency              float64   `db:"efficiency" json:"efficiency"`
	CapKills                int64     `db:"capKills" json:"capKills"`
	Kills                   int64     `db:"kills" json:"kills"`
	Losses                  int64     `db:"losses" json:"losses"`
}

// GetAlliance Obtain alliance information by ID.
func GetAlliance(id int64) (*Alliance, error) {
	ref := Alliance{}
	if err := database.QueryRowx(`
		SELECT 
			A.allianceID,
		    A.name AS allianceName, 
		    A.shortName AS allianceTicker,
		    A.corporationsCount,
			A.startDate,
			coalesce(efficiency, 0) AS efficiency,
			coalesce(capKills, 0) AS capKills,
			coalesce(kills, 0) AS kills,
			coalesce(losses, 0) AS losses,
		    
		    EXEC.name AS executorCorporationName,
		    EXEC.corporationID AS executorCorporationID
		    
		FROM evedata.alliances A
		INNER JOIN evedata.corporations EXEC ON A.executorCorpID = EXEC.corporationID
		LEFT OUTER JOIN evedata.entityKillStats S ON S.id = A.allianceID
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

type AllianceHistory struct {
	CorporationID   int64     `db:"corporationID" json:"corporationID"`
	CorporationName string    `db:"corporationName" json:"corporationName"`
	StartDate       time.Time `db:"startDate" json:"startDate"`
	EndDate         null.Time `db:"endDate" json:"endDate,omitempty"`
}

// Obtain a list of corporations history with an alliance by ID.

func GetAllianceHistory(id int64) ([]AllianceHistory, error) {
	ref := []AllianceHistory{}
	if err := database.Select(&ref, `
		SELECT H.corporationID, name AS corporationName, startDate, endDate
		FROM evedata.allianceHistory H
		INNER JOIN evedata.corporations C ON C.corporationID = H.corporationID
		WHERE H.allianceID = ?
		ORDER BY startDate DESC;
		`, id); err != nil {
		return nil, err
	}

	return ref, nil
}

type AllianceJoinHistory struct {
	CharacterID   int64     `db:"characterID" json:"characterID"`
	CharacterName string    `db:"characterName" json:"characterName"`
	Date          time.Time `db:"date" json:"date"`
	Event         string    `db:"event" json:"event"`
}

// Obtain a list of corporations history with an alliance by ID.

func GetAllianceJoinHistory(id int64) ([]AllianceJoinHistory, error) {
	ref := []AllianceJoinHistory{}
	if err := database.Select(&ref, `
		SELECT name AS characterName, C.characterID, event, date FROM (
			SELECT H.characterID, "joined" AS event, H.startDate AS date
					FROM evedata.corporationHistory H
					INNER JOIN evedata.allianceHistory A ON A.corporationID = H.corporationID
						AND H.startDate > A.startDate AND H.startDate < IFNULL(A.endDate, UTC_TIMESTAMP())
					WHERE A.allianceID = ?
			UNION		
			SELECT H.characterID, "left" AS event, H.endDate AS date
					FROM evedata.corporationHistory H
					INNER JOIN evedata.allianceHistory A ON A.corporationID = H.corporationID
						AND H.endDate > A.startDate AND H.endDate < IFNULL(A.endDate, UTC_TIMESTAMP())
					WHERE A.allianceID = ? AND H.endDate IS NOT NULL) S 
			INNER JOIN evedata.characters C ON S.characterID = C.characterID 
			ORDER BY date DESC
			
		`, id, id); err != nil {
		return nil, err
	}

	return ref, nil
}
