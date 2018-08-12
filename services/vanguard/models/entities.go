package models

type EntityName struct {
	Name       string `db:"name" json:"name"`
	EntityType string `db:"type" json:"type"`
}

// Obtain entity name and type by ID.

func GetEntityName(id int64) (*EntityName, error) {
	ref := EntityName{}
	if err := database.QueryRowx(`
		SELECT name, 'corporation' AS type FROM evedata.corporations WHERE corporationID = ?
		UNION
		SELECT name, 'alliance' AS type FROM evedata.alliances WHERE allianceID = ?
		LIMIT 1`, id, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

// Obtain type name.

func GetTypeName(id int64) (string, error) {
	ref := ""
	if err := database.QueryRowx(`
		SELECT typeName FROM invTypes WHERE typeID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}

// Obtain SolarSystem name.

func GetSystemName(id int64) (string, error) {
	ref := ""
	if err := database.QueryRowx(`
		SELECT solarSystemName FROM mapSolarSystems WHERE solarSystemID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}

// Obtain Celestial name.

func GetCelestialName(id int64) (string, error) {
	ref := ""
	if err := database.QueryRowx(`
		SELECT itemName FROM mapDenormalize WHERE itemID = ?
		LIMIT 1`, id).Scan(&ref); err != nil {
		return "", err
	}
	return ref, nil
}
