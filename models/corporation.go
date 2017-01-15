package models

import (
	"time"

	"github.com/guregu/null"
)

func UpdateCorporation(corporationID int64, name string, ticker string, ceoID int64, stationID int64,
	description string, allianceID int64, factionID int64, url string, memberCount int64, shares int64, cacheUntil time.Time) error {

	cacheUntil = time.Now().UTC().Add(time.Hour * 24 * 5)
	if _, err := database.Exec(`
		INSERT INTO evedata.corporations
			(corporationID,name,ticker,ceoID,stationID,description,allianceID,factionID,url,memberCount,shares,updated,cacheUntil)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
			ceoID=VALUES(ceoID), stationID=VALUES(stationID), description=VALUES(description), allianceID=VALUES(allianceID), 
			factionID=VALUES(factionID), url=VALUES(url), memberCount=VALUES(memberCount),  
			shares=VALUES(shares), updated=UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, corporationID, name, ticker, ceoID, stationID, description, allianceID, factionID, url, memberCount, shares, cacheUntil); err != nil {
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
	Description     string      `db:"description" json:"description"`
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
		    Al.name AS allianceName,
		    C.description
		FROM evedata.corporations C
		LEFT OUTER JOIN evedata.alliances Al ON C.allianceID = Al.allianceID
        LEFT OUTER JOIN evedata.characters Ch ON Ch.characterID = C.ceoID
		WHERE C.corporationID = ?
		LIMIT 1`, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}
