package models

import (
	"fmt"
	"time"
)

func GetKnownKillmails() ([]int64, error) {
	var known []int64
	if err := database.Select(&known, `SELECT id FROM evedata.killmails;`); err != nil {
		return nil, err
	}
	return known, nil
}

// KillmailList is brief list of killmails
type KillmailList struct {
	ID            int32     `db:"id" json:"id"`
	SolarSystemID int32     `db:"solarSystemID" json:"solarSystemID"`
	ShipType      int32     `db:"shipType" json:"shipType"`
	WarID         int32     `db:"warID" json:"warID"`
	IsLoss        int32     `db:"isLoss" json:"isLoss"`
	Killtime      time.Time `db:"killtime" json:"killtime"`
}

// GetKillmailsForEntity fetches all the killmails for an entity
func GetKillmailsForEntity(id int64, entityType string) ([]KillmailList, error) {
	kill := []KillmailList{}

	var victim, entity string

	switch entityType {
	case "corporation":
		victim = fmt.Sprintf("victimCorporationID=%d", id)
		entity = fmt.Sprintf("corporationID=%d", id)
	case "alliance":
		victim = fmt.Sprintf("victimAllianceID=%d", id)
		entity = fmt.Sprintf("allianceID=%d", id)
	case "character":
		victim = fmt.Sprintf("victimCharacterID=%d", id)
		entity = fmt.Sprintf("characterID=%d", id)
	}

	if err := database.Select(&kill, `
	SELECT * FROM (SELECT K.id, K.killtime, K.shipType, K.solarSystemID, K.warID, 1 AS isLoss
		FROM evedata.killmails K
		LEFT OUTER JOIN evedata.killmailAttackers A ON K.id = A.id
		WHERE `+victim+` 
		UNION DISTINCT
		SELECT DISTINCT K.id, K.killtime, K.shipType, K.solarSystemID, K.warID, 0 AS isLoss
		FROM evedata.killmails K
		INNER JOIN evedata.killmailAttackers A ON K.id = A.id
		WHERE `+entity+` 
		) a
		ORDER BY killTime DESC`); err != nil {
		return nil, err
	}

	return kill, nil
}

// [TODO] Break out the CSV into an array
type KillmailHeatMap struct {
	Day   int64 `db:"day" json:"day"`
	Hour  int64 `db:"hour" json:"hour"`
	Value int64 `db:"value" json:"value"`
}

func GetKillmailHeatMap(id int64, entityType string) ([]KillmailHeatMap, error) {
	v := []KillmailHeatMap{}

	var victim, entity string

	switch entityType {
	case "corporation":
		victim = fmt.Sprintf("victimCorporationID=%d", id)
		entity = fmt.Sprintf("corporationID=%d", id)
	case "alliance":
		victim = fmt.Sprintf("victimAllianceID=%d", id)
		entity = fmt.Sprintf("allianceID=%d", id)
	case "character":
		victim = fmt.Sprintf("victimCharacterID=%d", id)
		entity = fmt.Sprintf("characterID=%d", id)
	}

	if err := database.Select(&v, `
		SELECT DAYOFWEEK(killTime) AS day, HOUR(killTime) AS hour, count(S.id) AS value FROM (
			SELECT killTime, K.id
			FROM evedata.killmails K
			WHERE `+victim+`
			UNION
			SELECT DISTINCT killTime, K.id
			FROM evedata.killmails K
			INNER JOIN evedata.killmailAttackers A ON K.id = A.id
			WHERE `+entity+`
		) S GROUP BY DAYOFWEEK(killTime), HOUR(killTime)`); err != nil {
		return nil, err
	}
	return v, nil
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
