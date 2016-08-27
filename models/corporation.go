package models

import "time"

func UpdateCorporation(corporationID int64, name string, ticker string, ceoID int64, stationID int64,
	description string, allianceID int64, factionID int64, url string, memberCount int64, shares int64, cacheUntil time.Time) error {
	if _, err := database.Exec(`
		INSERT INTO corporation
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
