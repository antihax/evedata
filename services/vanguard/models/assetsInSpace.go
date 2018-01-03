package models

import (
	"time"

	"github.com/guregu/null"
)

type AssetsInSpace struct {
	ItemName string    `db:"itemName" json:"itemName"`
	TypeName string    `db:"typeName" json:"typeName"`
	TypeID   int64     `db:"typeID" json:"typeID"`
	Security float64   `db:"security" json:"security"`
	LastSeen time.Time `db:"lastSeen" json:"lastSeen"`
}

// Obtain a list of alliance assets in space.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetAllianceAssetsInSpace(id int64) ([]AssetsInSpace, error) {
	ref := []AssetsInSpace{}
	if err := database.Select(&ref, `
		SELECT itemName, A.typeID, typeName, lastSeen, round(security,1) AS security
			FROM evedata.discoveredAssets A
		    INNER JOIN invTypes T ON A.typeID = T.typeID
			INNER JOIN mapDenormalize D ON A.locationID = D.itemID
			WHERE allianceID = ?
		`, id); err != nil {
		return nil, err
	}
	return ref, nil
}

// Obtain a list of corporation assets in space.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCorporationAssetsInSpace(id int64) ([]AssetsInSpace, error) {
	ref := []AssetsInSpace{}
	if err := database.Select(&ref, `
		SELECT itemName, A.typeID, typeName, lastSeen, round(security,1) AS security
			FROM evedata.discoveredAssets A
		    INNER JOIN invTypes T ON A.typeID = T.typeID
			INNER JOIN mapDenormalize D ON A.locationID = D.itemID
			WHERE corporationID = ?
		`, id); err != nil {
		return nil, err
	}
	return ref, nil
}

type LostFighters struct {
	KillID          int64       `db:"killID" json:"killID"`
	TypeName        string      `db:"typeName" json:"typeName"`
	AllianceName    null.String `db:"allianceName" json:"allianceName"`
	AllianceID      null.Int    `db:"allianceID" json:"allianceID"`
	CorporationName string      `db:"corporationName" json:"corporationName"`
	CorporationID   int64       `db:"corporationID" json:"corporationID"`
	SolarSystem     string      `db:"solarSystemName" json:"solarSystemName"`
	KillTime        time.Time   `db:"killTime" json:"killTime"`
}

// Obtain a list lost fighters in highsec
// [BENCHMARK] 0.437 sec / 0.000 sec
func GetCorporationAssetsInSpaceLostFightersHighsec() ([]LostFighters, error) {
	ref := []LostFighters{}
	if err := database.Select(&ref, `
		SELECT K.id AS killID, typeName, A.name AS allianceName, A.allianceID, C.name AS corporationName, corporationID, solarSystemName, killTime FROM evedata.killmails K
			INNER JOIN eve.mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			LEFT OUTER JOIN evedata.alliances A ON K.victimAllianceID = A.allianceID
			LEFT OUTER JOIN evedata.corporations C ON K.victimCorporationID = C.corporationID
			INNER JOIN eve.invTypes T ON T.typeID = K.shipType
		WHERE round(security,1) >= 0.5 AND 
			victimCharacterID = 0 AND 
			groupID IN (1537, 1652, 1653) AND 
			killTime > DATE_SUB(NOW(), INTERVAL 8 DAY)
		UNION
		SELECT DISTINCT K.id AS killID, typeName, A.name AS allianceName, A.allianceID, C.name AS corporationName, C.corporationID, solarSystemName, killTime FROM evedata.killmails K
			INNER JOIN evedata.killmailAttackers KA ON K.id = KA.id AND KA.characterID = 0
			INNER JOIN eve.mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			LEFT OUTER JOIN evedata.alliances A ON KA.allianceID = A.allianceID
			LEFT OUTER JOIN evedata.corporations C ON KA.corporationID = C.corporationID
			INNER JOIN eve.invTypes T ON T.typeID = KA.shipType AND T.groupID IN (365, 549, 1023, 1537, 1652, 1653, 1657, 2233)
		WHERE round(security,1) >= 0.5 AND 
			killTime > DATE_SUB(NOW(), INTERVAL 8 DAY)
		ORDER BY killTime DESC
		`); err != nil {
		return nil, err
	}
	return ref, nil
}
