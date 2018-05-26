package models

type SystemVertex struct {
	FromSolarSystemID int    `db:"fromSolarSystemID"`
	Connections       string `db:"connections"`
}

// [BENCHMARK] 0.407 sec / 0.421 sec [TODO] Optimize
func GetSystemVertices() ([]SystemVertex, error) {
	s := []SystemVertex{}
	if err := database.Select(&s, `
		SELECT fromSolarSystemID, GROUP_CONCAT(toSolarSystemID) AS connections
		FROM eve.mapSolarSystemJumps
		GROUP BY fromSolarSystemID
	`); err != nil {
		return nil, err
	}
	return s, nil
}

type SystemCelestials struct {
	Type string  `db:"type" json:"type"`
	Name string  `db:"name" json:"name"`
	X    float64 `db:"x" json:"x"`
	Y    float64 `db:"y" json:"y"`
	Z    float64 `db:"z" json:"z"`
}

// [BENCHMARK] 0.407 sec / 0.421 sec [TODO] Optimize
func GetSystemCelestials(system int32) ([]SystemCelestials, error) {
	s := []SystemCelestials{}
	if err := database.Select(&s, `
		SELECT groupName AS type, IFNULL(D.itemName, solarSystemName) AS name, D.x, D.y, D.z FROM mapDenormalize D
		INNER JOIN invGroups G ON G.groupID = D.groupID
		LEFT OUTER JOIN mapJumps J ON J.starGateID = D.itemID
		LEFT OUTER JOIN mapDenormalize D2 ON D2.itemID = J.destinationID
		LEFT OUTER JOIN mapSolarSystems S ON S.solarSystemID = D2.solarSystemID
		WHERE D.solarSystemID = ? 
	`, system); err != nil {
		return nil, err
	}
	return s, nil
}

type NullSystems struct {
	SolarSystemID   int32  `db:"solarSystemID" json:"solarSystemID"`
	SolarSystemName string `db:"solarSystemName" json:"solarSystemName"`
	RegionName      string `db:"regionName" json:"regionName"`
}

// [BENCHMARK] 0.407 sec / 0.421 sec [TODO] Optimize
func GetNullSystems() ([]NullSystems, error) {
	s := []NullSystems{}
	if err := database.Select(&s, `
		SELECT solarSystemID, solarSystemName, regionName
		FROM mapSolarSystems S
		INNER JOIN mapRegions R ON R.regionID = S.regionID
		WHERE security < 0 AND solarSystemID < 31000000
		ORDER BY solarSystemName
	`); err != nil {
		return nil, err
	}
	return s, nil
}
