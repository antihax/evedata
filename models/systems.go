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
