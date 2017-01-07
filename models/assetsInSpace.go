package models

import "time"

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
