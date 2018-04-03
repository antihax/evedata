package models

import (
	"time"

	"github.com/guregu/null"
)

func UpdateCorporation(corporationID int32, name string, ticker string, ceoID int32,
	allianceID int32, factionID int32, memberCount int32, cacheUntil time.Time) error {
	cacheUntil = time.Now().UTC().Add(time.Hour * 24 * 5)
	if _, err := database.Exec(`
		INSERT INTO evedata.corporations
			(corporationID,name,ticker,ceoID,allianceID,factionID,memberCount,updated,cacheUntil)
			VALUES(?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
			ceoID=VALUES(ceoID), name=VALUES(name), ticker=VALUES(ticker), allianceID=VALUES(allianceID), 
			factionID=VALUES(factionID), memberCount=VALUES(memberCount),  
			updated=UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, corporationID, name, ticker, ceoID, allianceID, factionID, memberCount, cacheUntil); err != nil {
		return err
	}
	return nil
}

type Corporation struct {
	CorporationID   int64       `db:"corporationID" json:"corporationID"`
	CorporationName string      `db:"corporationName" json:"corporationName"`
	AllianceID      int64       `db:"allianceID" json:"allianceID"`
	AllianceName    null.String `db:"allianceName" json:"allianceName"`
	CEOID           int64       `db:"ceoID" json:"ceoID"`
	CEOName         string      `db:"ceoName" json:"ceoName"`
	MemberCount     int64       `db:"memberCount" json:"memberCount"`
}

// Obtain Corporation information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCorporation(id int64) (*Corporation, error) {
	ref := Corporation{}
	if err := database.QueryRowx(`
		SELECT 
			C.corporationID,
		    C.name AS corporationName,
		    C.memberCount,
            IFNULL(ceoID,0) AS ceoID,
            IFNULL(Ch.name, "") AS ceoName,
		    IFNULL(Al.allianceID,0) AS allianceID,
		    Al.name AS allianceName
		FROM evedata.corporations C
		LEFT OUTER JOIN evedata.alliances Al ON C.allianceID = Al.allianceID
        LEFT OUTER JOIN evedata.characters Ch ON Ch.characterID = C.ceoID
		WHERE C.corporationID = ?
		LIMIT 1`, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

type CorporationJoinHistory struct {
	CharacterID   int64     `db:"characterID" json:"characterID"`
	CharacterName string    `db:"characterName" json:"characterName"`
	StartDate     time.Time `db:"startDate" json:"startDate"`
	EndDate       null.Time `db:"endDate" json:"endDate,omitempty"`
}

// Obtain a list of corporations history with an alliance by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCorporationJoinHistory(id int64) ([]CorporationJoinHistory, error) {
	ref := []CorporationJoinHistory{}
	if err := database.Select(&ref, `
		SELECT H.characterID, name AS characterName, startDate, endDate
		FROM evedata.corporationHistory H
		INNER JOIN evedata.characters C ON C.characterID = H.characterID
		WHERE H.corporationID = ?
		ORDER BY startDate DESC;
		`, id); err != nil {
		return nil, err
	}

	return ref, nil
}
