package models

import "time"

func UpdateAlliance(allianceID int64, name string, memberCount int64, shortName string, executorCorp int64, startDate time.Time) error {
	if _, err := database.Exec(`
		INSERT INTO alliances 
			(allianceID,name,shortName,executorCorpID,memberCount,startDate,updated)
			VALUES(?,?,?,?,?,?,UTC_TIMESTAMP()) 
			ON DUPLICATE KEY UPDATE 
			executorCorpID = VALUE(executorCorpID), memberCount = VALUE(memberCount), updated = UTC_TIMESTAMP()
	`, allianceID, name, shortName, executorCorp, startDate); err != nil {
		return err
	}
	return nil
}
