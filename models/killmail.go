package models

import "time"

func AddKillmail(id int64, solarSystemID int64, killTime time.Time, victimCharacterID int64, victimCorporationID int64,
	victimAllianceID int64, hash string, attackerCount int64, damageTaken int64, x float64, y float64, z float64,
	shipType int64, warID int64) error {
	if _, err := database.Exec(`
		INSERT INTO killmails
			(id,solarSystemID,killTime,victimCharacterID,victimCorporationID,victimAllianceID,hash,
			attackerCount,damageTaken,x,y,z,shipType,warID)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`, id, solarSystemID, killTime, victimCharacterID, victimCorporationID, victimAllianceID, hash,
		attackerCount, damageTaken, x, y, z, shipType, warID); err != nil {
		return err
	}
	return nil
}

func AddKillmailAttacker(id int64, characterID int64, corporationID int64, allianceID int64, shipType int64,
	finalBlow bool, damageDone int64, weaponType int64, securityStatus float64) error {
	if _, err := database.Exec(`
		INSERT INTO killmailAttackers
			(id,characterID,corporationID,allianceID,shipType,finalBlow,damageDone,weaponType,securityStatus)
			VALUES(?,?,?,?,?,?,?,?,?);
	`, id, characterID, corporationID, allianceID, shipType, finalBlow, damageDone, weaponType, securityStatus); err != nil {
		return err
	}
	return nil
}

func AddKillmailItems(id int64, itemType int64, flag int64, quantityDestroyed int64, quantityDropped int64, singleton int64) error {
	if _, err := database.Exec(`
		INSERT INTO killmailItems
			(id,itemType,flag,quantityDestroyed,quantityDropped,singleton)
			VALUES(?,?,?,?,?,?);	
	`, id, itemType, flag, quantityDestroyed, quantityDropped, singleton); err != nil {
		return err
	}
	return nil
}

func GetKnownKillmails() ([]int64, error) {
	var known []int64
	if err := database.Select(&known, `SELECT id FROM killmails;`); err != nil {
		return nil, err
	}
	return known, nil
}
