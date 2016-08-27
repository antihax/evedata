package models

import "time"

func UpdateAlliance(allianceID int64, name string, memberCount int64, shortName string, executorCorp int64, startDate time.Time, deleted bool,
	description string, creatorCorp int64, creatorCharacter int64, cacheUntil time.Time) error {
	if _, err := database.Exec(`
		INSERT INTO alliance 
			(allianceID,name,shortName,executorCorpID,startDate,corporationsCount,deleted,description,creatorCorpID,creatorCharacter,updated,cacheUntil)
			VALUES(?,?,?,?,?,?,?,?,?,?,UTC_TIMESTAMP(),?) 
			ON DUPLICATE KEY UPDATE 
			executorCorpID = VALUES(executorCorpID), corporationsCount = VALUES(corporationsCount), 
			description = VALUES(description), deleted = VALUES(deleted), updated = UTC_TIMESTAMP(), cacheUntil=VALUES(cacheUntil)
	`, allianceID, name, shortName, executorCorp, startDate, memberCount, deleted, description, creatorCorp, creatorCharacter, cacheUntil); err != nil {
		return err
	}
	return nil
}
