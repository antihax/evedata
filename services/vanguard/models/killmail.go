package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/guregu/null"
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

// KillmailDetails is all the details of a killmail
type KillmailDetails struct {
	SolarSystemID   int32          `db:"solarSystemID" json:"solarSystemID"`
	SolarSystemName string         `db:"solarSystemName" json:"solarSystemName"`
	Security        float32        `db:"security" json:"security"`
	TypeID          int32          `db:"typeID" json:"typeID"`
	TypeName        string         `db:"typeName" json:"typeName"`
	CharacterID     sql.NullInt64  `db:"characterID" json:"characterID"`
	CharacterName   sql.NullString `db:"characterName" json:"characterName"`
	CorporationID   sql.NullInt64  `db:"corporationID" json:"corporationID"`
	CorporationName sql.NullString `db:"corporationName" json:"corporationName"`
	AllianceID      sql.NullInt64  `db:"allianceID" json:"allianceID"`
	AllianceName    sql.NullString `db:"allianceName" json:"allianceName"`
	FactionID       sql.NullInt64  `db:"factionID" json:"factionID"`
	FactionName     sql.NullString `db:"factionName" json:"factionName"`
	Hash            string         `db:"hash" json:"hash"`
}

// GetKillmailDetails fetches all the details of a killmail
func GetKillmailDetails(id int64) (*KillmailDetails, error) {
	kill := KillmailDetails{}
	if err := database.QueryRowx(
		`	SELECT 
				Sy.solarSystemID, Sy.solarSystemName,
				ROUND(Sy.security,1) AS security,
				S.typeID, S.typeName, 
				C.characterID, C.name AS characterName,
				Co.corporationID, Co.name AS corporationName,
				A.allianceID, A.name AS allianceName,
				K.factionID, itemName AS factionName,
				K.hash
				FROM evedata.killmails K
				LEFT OUTER JOIN evedata.characters C ON K.victimCharacterID = C.characterID
				LEFT OUTER JOIN evedata.corporations Co ON K.victimCorporationID = Co.corporationID
				LEFT OUTER JOIN evedata.alliances A ON K.victimAllianceID = A.allianceID
				LEFT OUTER JOIN eve.eveNames N ON K.factionID = N.itemID
				INNER JOIN eve.invTypes S ON S.typeID = K.shipType
				INNER JOIN eve.mapSolarSystems Sy ON Sy.solarSystemID = K.solarSystemID
				WHERE id = ? LIMIT 1`, id).StructScan(&kill); err != nil {

		return nil, err
	}
	return &kill, nil
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
				AND round(S.security,1) >= 0.5 AND K.victimAllianceID > 0
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
