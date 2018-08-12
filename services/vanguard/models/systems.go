package models

type SystemVertex struct {
	FromSolarSystemID int    `db:"fromSolarSystemID"`
	Connections       string `db:"connections"`
}

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
	DNA  string  `db:"DNA" json:"DNA"`
	X    float64 `db:"x" json:"x"`
	Y    float64 `db:"y" json:"y"`
	Z    float64 `db:"z" json:"z"`
	SX   float64 `db:"sx" json:"sx,omitempty"`
	SY   float64 `db:"sy" json:"sy,omitempty"`
	SZ   float64 `db:"sz" json:"sz,omitempty"`
	DX   float64 `db:"dx" json:"dx,omitempty"`
	DY   float64 `db:"dy" json:"dy,omitempty"`
	DZ   float64 `db:"dz" json:"dz,omitempty"`
}

func GetSystemCelestials(system int32) ([]SystemCelestials, error) {
	s := []SystemCelestials{}
	if err := database.Select(&s, `
		SELECT coalesce(IF(sofHullName != "", concat(sofHullName, ":",sofFactionName, ":", sofRaceName), NULL),  graphicFile, "") AS DNA, 
		groupName AS type, IFNULL(D.itemName, S.solarSystemName) AS name, 
		D.x, D.y, D.z,
		coalesce(O.x,0) AS sx,   
		coalesce(O.y,0) AS sy,
		coalesce(O.z,0) AS sz,
		coalesce(S.x,0) AS dx,   
		coalesce(S.y,0) AS dy,
		coalesce(S.z,0) AS dz
			FROM mapDenormalize D
			INNER JOIN invGroups G ON G.groupID = D.groupID
			LEFT OUTER JOIN mapJumps J ON J.starGateID = D.itemID
			LEFT OUTER JOIN mapDenormalize D2 ON D2.itemID = J.destinationID
			LEFT OUTER JOIN mapSolarSystems S ON S.solarSystemID = D2.solarSystemID
			LEFT OUTER JOIN mapSolarSystems O ON O.solarSystemID = D.solarSystemID
			LEFT OUTER JOIN invTypes T ON T.typeID = D.typeID
			LEFT OUTER JOIN eveGraphics GR ON T.graphicID = GR.graphicID
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
