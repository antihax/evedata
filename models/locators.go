package models

import (
	"time"
)

type LocatorShares struct {
	CharacterID int64  `db:"characterID" json:"characterID,omitempty"`
	ID          int64  `db:"entityID" json:"id,omitempty"`
	EntityName  string `db:"entityName" json:"entityName,omitempty"`
	EntityType  string `db:"type" json:"type,omitempty"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetLocatorShares(characterID int64) ([]LocatorShares, error) {
	locatorShares := []LocatorShares{}
	if err := database.Select(&locatorShares, `
		SELECT characterID, entityID, IFNULL(A.name, C.name) AS entityName, IF(A.name IS NULL, "corporation", "alliance") AS type
		FROM evedata.locatorShareWith S
		LEFT OUTER JOIN evedata.corporations C ON C.corporationID = S.entityID
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = S.entityID
		WHERE characterID = ?;`, characterID); err != nil {
		return nil, err
	}
	return locatorShares, nil
}

type LocatorResults struct {
	SystemID          int64     `db:"systemID" json:"systemID,omitempty"`
	SystemName        string    `db:"systemName" json:"systemName,omitempty"`
	RegionID          int64     `db:"regionID" json:"regionID,omitempty"`
	RegionName        string    `db:"regionName" json:"regionName,omitempty"`
	ConstellationID   int64     `db:"constellationID" json:"constellationID,omitempty"`
	ConstellationName string    `db:"constellationName" json:"constellationName,omitempty"`
	CharacterID       int64     `db:"characterID" json:"characterID,omitempty"`
	CharacterName     string    `db:"characterName" json:"characterName,omitempty"`
	CorporationID     int64     `db:"corporationID" json:"corporationID,omitempty"`
	CorporationName   string    `db:"corporationName" json:"corporationName,omitempty"`
	AllianceID        int64     `db:"allianceID" json:"allianceID,omitempty"`
	AllianceName      string    `db:"allianceName" json:"allianceName,omitempty"`
	Time              time.Time `db:"time" json:"time"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetLocatorResults(characterID int64) ([]LocatorResults, error) {
	locatorResults := []LocatorResults{}
	if err := database.Select(&locatorResults, `
		SELECT DISTINCT Sy.solarSystemID AS systemID, Sy.solarSystemName AS systemName, R.regionID, R.regionName, 
						Con.constellationID, Con.constellationName, Ca.name AS characterName, L.characterID, 
						Co.name AS corporationName, Ca.corporationID, A.name AS allianceName, Co.allianceID, time 
		FROM evedata.locatedCharacters L
		INNER JOIN evedata.locatorShareWith S ON L.characterID = S.characterID 
		INNER JOIN evedata.characters Ca ON L.locatedCharacterID = Ca.characterID
		INNER JOIN evedata.corporations Co ON Ca.corporationID = Co.corporationID
		LEFT OUTER JOIN evedata.alliances A ON A.allianceID = Co.allianceID
		INNER JOIN eve.mapSolarSystems Sy ON Sy.solarSystemID = L.solarSystemID
		INNER JOIN eve.mapConstellations Con ON Con.constellationID = L.constellationID
		INNER JOIN eve.mapRegions R ON R.regionID = L.regionID
		WHERE L.characterID = ? OR S.entityID IN 
		(SELECT corporationID  FROM evedata.characters WHERE characterID = ?
		UNION
		SELECT allianceID FROM evedata.characters WHERE characterID = ?);`, characterID, characterID, characterID); err != nil {
		return nil, err
	}
	return locatorResults, nil
}

func AddLocatorShare(characterID int64, entityID int64) error {
	if _, err := database.Exec(`
		INSERT INTO evedata.locatorShareWith	(characterID, entityID)
			VALUES(?,?)
			ON DUPLICATE KEY UPDATE entityID = entityID`,
		characterID, entityID); err != nil {
		return err
	}
	return nil
}

func DeleteLocatorShare(characterID int64, entityID int64) error {
	if _, err := database.Exec(`DELETE FROM evedata.locatorShareWith WHERE characterID = ? AND entityID = ? LIMIT 1`,
		characterID, entityID); err != nil {
		return err
	}
	return nil
}
