package models

import (
	"evedata/null"
	"time"
)

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

type ActiveWarList struct {
	WarID         int64       `db:"warID" json:"warID"`
	TimeStarted   time.Time   `db:"timeStarted" json:"timeStarted"`
	TimeFinished  time.Time   `db:"timeFinished" json:"timeFinished"`
	OpenForAllies bool        `db:"openForAllies" json:"openForAllies"`
	AggressorID   int64       `db:"aggressorID" json:"aggressorID"`
	AggressorType null.String `db:"aggressorType" json:"aggressorType"`
	AggressorName null.String `db:"aggressorName" json:"aggressorName"`
	DefenderID    int64       `db:"defenderID" json:"defenderID"`
	DefenderType  null.String `db:"defenderType" json:"defenderType"`
	DefenderName  null.String `db:"defenderName" json:"defenderName"`
	Mutual        bool        `db:"mutual" json:"mutual"`
	Kills         int64       `db:"kills" json:"kills"`
	Losses        int64       `db:"losses" json:"losses"`
}

func GetActiveWarList() ([]ActiveWarList, error) {

	wars := []ActiveWarList{}
	if err := database.Select(&wars, `
	SELECT 
		W.id AS warID, 
	    timeStarted, 
	    timeFinished, 
	    openForAllies, 
	    aggressorID, 
	    Ag.Type AS aggressorType, 
	    defenderID, 
	    Df.type AS defenderType, 
	    mutual, 
	    IFNULL(kills,0) as kills,  
	    IFNULL(losses,0) as losses,
	    IF(AA.allianceID > 0, AA.name, AC.name) AS aggressorName,
	    IF(DA.allianceID > 0, DA.name, DC.name) AS defenderName
	        
		FROM wars W
		INNER JOIN crestID Ag ON Ag.id = aggressorID
	    INNER JOIN crestID Df ON Df.id = defenderID
	    LEFT OUTER JOIN alliance AA on AA.allianceID = aggressorID
		LEFT OUTER JOIN alliance DA on DA.allianceID = defenderID
		LEFT OUTER JOIN corporation AC on AC.corporationID = aggressorID
		LEFT OUTER JOIN corporation DC on DC.corporationID = defenderID
		LEFT OUTER JOIN 
	    ( -- Kills by the Aggressor
			SELECT 
				W.id, 
				count(*) AS kills
				FROM wars W
				INNER JOIN killmails K ON K.warID = W.id AND 
				(
					K.victimAllianceID != W.aggressorID AND 
					K.victimCorporationID != W.aggressorID
				)
				GROUP BY W.id
		) AS K ON W.id = K.id
		LEFT OUTER JOIN 
	    ( -- Kills by the Defenders
			SELECT 
				W.id, 
				count(*) AS losses
				FROM wars W
				INNER JOIN killmails L ON L.warID = W.id AND 
				(
					L.victimAllianceID = W.aggressorID OR 
					L.victimCorporationID = W.aggressorID
				)
				GROUP BY W.id
		) AS L ON W.id = L.id
	    WHERE 
			timeFinished > UTC_TIMESTAMP() OR
	        timeFinished = "0001-01-01 00:00:00" AND
	        mutual = 0`); err != nil {
		return nil, err
	}
	return wars, nil
}
