package models

import (
	"fmt"
	"time"

	"github.com/guregu/null"
)

func AddKillmail(id int32, solarSystemID int32, killTime time.Time, victimCharacterID int32, victimCorporationID int32,
	victimAllianceID int32, hash string, attackerCount int, damageTaken int32, x float32, y float32, z float32,
	shipType int32, warID int32) error {
	if _, err := database.Exec(`
		INSERT IGNORE INTO evedata.killmails
			(id,solarSystemID,killTime,victimCharacterID,victimCorporationID,victimAllianceID,hash,
			attackerCount,damageTaken,x,y,z,shipType,warID)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`, id, solarSystemID, killTime, victimCharacterID, victimCorporationID, victimAllianceID, hash,
		attackerCount, damageTaken, x, y, z, shipType, warID); err != nil {
		return err
	}
	return nil
}

func AddKillmailAttacker(id int32, characterID int32, corporationID int32, allianceID int32, shipType int32,
	finalBlow bool, damageDone int32, weaponType int32, securityStatus float32) error {
	if _, err := database.Exec(`
		INSERT IGNORE INTO evedata.killmailAttackers
			(id,characterID,corporationID,allianceID,shipType,finalBlow,damageDone,weaponType,securityStatus)
			VALUES(?,?,?,?,?,?,?,?,?);
	`, id, characterID, corporationID, allianceID, shipType, finalBlow, damageDone, weaponType, securityStatus); err != nil {
		return err
	}
	return nil
}

func AddKillmailItems(id int32, itemType int32, flag int32, quantityDestroyed int64, quantityDropped int64, singleton int32) error {
	if _, err := database.Exec(`
		INSERT IGNORE INTO evedata.killmailItems
			(id,itemType,flag,quantityDestroyed,quantityDropped,singleton)
			VALUES(?,?,?,?,?,?);	
	`, id, itemType, flag, quantityDestroyed, quantityDropped, singleton); err != nil {
		return err
	}
	return nil
}

func GetKnownKillmails() ([]int64, error) {
	var known []int64
	if err := database.Select(&known, `SELECT id FROM evedata.killmails;`); err != nil {
		return nil, err
	}
	return known, nil
}

// [TODO] Break out the CSV into an array
type ConstellationActivity struct {
	Number            int64  `db:"number" json:"number"`
	SolarSystemIDs    string `db:"solarSystemIDs" json:"solarSystemIDs"`
	SolarSystemNames  string `db:"solarSystemNames" json:"solarSystemNames"`
	ConstellationID   int64  `db:"constellationID" json:"constellationID"`
	ConstellationName string `db:"constellationName" json:"constellationName"`
	RegionID          int64  `db:"regionID" json:"regionID"`
	RegionName        string `db:"regionName" json:"regionName"`
}

// [BENCHMARK] 0.015 sec / 0.000 sec
func GetConstellationActivity(id int64, entityType string) ([]ConstellationActivity, error) {
	v := []ConstellationActivity{}

	var victim, entity string

	switch entityType {
	case "corporation":
		victim = fmt.Sprintf("K.victimCorporationID=%d", id)
		entity = fmt.Sprintf("A.corporationID=%d", id)
	case "alliance":
		victim = fmt.Sprintf("K.victimAllianceID=%d", id)
		entity = fmt.Sprintf("A.allianceID=%d", id)
	case "character":
		victim = fmt.Sprintf("K.victimCharacterID=%d", id)
		entity = fmt.Sprintf("A.characterID=%d", id)
	}

	if err := database.Select(&v, `
	SELECT 	
		COUNT(DISTINCT K.id) AS number, 
	    GROUP_CONCAT(DISTINCT S.solarSystemID) AS solarSystemIDs,
	    GROUP_CONCAT(DISTINCT S.solarSystemName) AS solarSystemNames,
	    S.constellationID,
	    constellationName,
	    S.regionID,
	    regionName 
	FROM
	(SELECT K.id AS id, K.solarSystemID AS solarSystemID
		FROM evedata.killmails K 
		WHERE 
			K.killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 31 DAY) AND
			`+victim+`
	UNION ALL
		SELECT K.id, K.solarSystemID
			FROM evedata.killmails K 
			INNER JOIN evedata.killmailAttackers A ON A.id = K.id 
			WHERE 
				K.killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 31 DAY) AND
				`+entity+`
	) K
	INNER JOIN mapSolarSystems S ON K.solarSystemID = S.solarSystemID
	INNER JOIN mapConstellations C ON C.constellationID = S.constellationID
	INNER JOIN mapRegions R ON R.regionID = S.regionID
	GROUP BY S.constellationID`); err != nil {
		return nil, err
	}
	return v, nil
}

type KnownShipTypes struct {
	Number   int64  `db:"number" json:"number"`
	ShipType int64  `db:"shipType" json:"shipType"`
	ShipName string `db:"shipName" json:"shipName"`
}

// [BENCHMARK] 0.016 sec / 0.000 sec
func GetKnownShipTypes(id int64, entityType string) ([]KnownShipTypes, error) {
	v := []KnownShipTypes{}

	var victim, entity string

	switch entityType {
	case "corporation":
		victim = fmt.Sprintf("K.victimCorporationID=%d", id)
		entity = fmt.Sprintf("A.corporationID=%d", id)
	case "alliance":
		victim = fmt.Sprintf("K.victimAllianceID=%d", id)
		entity = fmt.Sprintf("A.allianceID=%d", id)
	case "character":
		victim = fmt.Sprintf("K.victimCharacterID=%d", id)
		entity = fmt.Sprintf("A.characterID=%d", id)
	}

	if err := database.Select(&v, `
	SELECT 	
		COUNT(DISTINCT K.id) AS number, 
	    shipType,
	    typeName AS shipName
	FROM
	(SELECT K.id AS id, K.shipType
		FROM evedata.killmails K 
		WHERE 
			K.killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 31 DAY) AND
			`+victim+`
			AND shipType > 0
	UNION ALL
		SELECT K.id, A.shipType
			FROM evedata.killmails K 
			INNER JOIN evedata.killmailAttackers A ON A.id = K.id 
			WHERE 
				K.killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 31 DAY) AND
				`+entity+`
                AND A.shipType > 0
	) K
	INNER JOIN invTypes T ON K.shipType = T.typeID
	GROUP BY shipType`); err != nil {
		return nil, err
	}
	return v, nil
}

type LossesInHighsec struct {
	Number     int64       `db:"number" json:"number"`
	ID         int64       `db:"id" json:"id"`
	Type       null.String `db:"type" json:"type"`
	Name       null.String `db:"name" json:"name"`
	Members    int64       `db:"members" json:"members"`
	Efficiency float64     `db:"efficiency" json:"efficiency"`
	Kills      int64       `db:"kills" json:"kills"`
	Losses     int64       `db:"losses" json:"losses"`
}

// [BENCHMARK] 0.10.703 sec / 0.000 sec
// FIXME: i am slow :(
func GetLossesInHighsec() ([]LossesInHighsec, error) {
	v := []LossesInHighsec{}

	if err := database.Select(&v, `
		SELECT 	
			COUNT(DISTINCT K.id) AS number, 
			A.name AS name,
			A.allianceID AS id,
			memberCount AS members,
			type,
			IFNULL(S.efficiency,1) AS efficiency,
			IFNULL(S.kills,0) AS kills,
			IFNULL(S.losses,0)  AS losses
		FROM
		(SELECT K.id AS id, victimAllianceID AS entityID, "alliance" AS type
			FROM evedata.killmails K
			INNER JOIN eve.mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE 
				K.killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY) 
				AND round(S.security,1) >= 0.5
		) K
		INNER JOIN evedata.alliances A ON K.entityID = A.allianceID
		LEFT OUTER JOIN evedata.entityKillStats S ON S.id = K.entityID 
		GROUP BY A.allianceID 
		HAVING members > 25 AND members < 500 AND number > 15
		UNION ALL

			SELECT 	
			COUNT(DISTINCT K.id) AS number, 
			A.name AS name,
			A.corporationID AS id,
			memberCount AS members,
			type,
			IFNULL(S.efficiency,1) AS efficiency,
			IFNULL(S.kills,0) AS kills,
			IFNULL(S.losses,0)  AS losses
		FROM
		(SELECT K.id AS id, victimCorporationID AS entityID, "corporation" AS type 
			FROM evedata.killmails K
			INNER JOIN eve.mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE 
				K.killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY) 
				AND round(S.security,1) >= 0.5
				AND victimAllianceID = 0
		) K
		INNER JOIN evedata.corporations A ON K.entityID = A.corporationID
		LEFT OUTER JOIN evedata.entityKillStats S ON S.id = K.entityID 
		WHERE corporationID > 2000000 
		GROUP BY A.corporationID 
		HAVING members > 25 AND members < 4000 AND number > 15
`); err != nil {
		return nil, err
	}
	return v, nil
}
