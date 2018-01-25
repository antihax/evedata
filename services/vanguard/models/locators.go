package models

import (
	"time"

	"github.com/guregu/null"
)

type LocatorResults struct {
	SystemID          int64       `db:"systemID" json:"systemID,omitempty"`
	SystemName        string      `db:"systemName" json:"systemName,omitempty"`
	RegionID          int64       `db:"regionID" json:"regionID,omitempty"`
	RegionName        string      `db:"regionName" json:"regionName,omitempty"`
	ConstellationID   int64       `db:"constellationID" json:"constellationID,omitempty"`
	ConstellationName string      `db:"constellationName" json:"constellationName,omitempty"`
	CharacterID       int32       `db:"characterID" json:"characterID,omitempty"`
	CharacterName     string      `db:"characterName" json:"characterName,omitempty"`
	CorporationID     int32       `db:"corporationID" json:"corporationID,omitempty"`
	CorporationName   string      `db:"corporationName" json:"corporationName,omitempty"`
	AllianceID        int32       `db:"allianceID" json:"allianceID,omitempty"`
	AllianceName      null.String `db:"allianceName" json:"allianceName,omitempty"`
	Time              time.Time   `db:"time" json:"time"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetLocatorResponses(characterID int32, cursorCharacterID int32) ([]LocatorResults, error) {
	locatorResults := []LocatorResults{}
	if err := database.Select(&locatorResults, `
		SELECT DISTINCT Sy.solarSystemID AS systemID, Sy.solarSystemName AS systemName, R.regionID, R.regionName, 
						Con.constellationID, Con.constellationName, Ca.name AS characterName, Ca.characterID, 
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
		SELECT allianceID FROM evedata.characters WHERE characterID = ?);`, characterID, cursorCharacterID, cursorCharacterID); err != nil {
		return nil, err
	}
	return locatorResults, nil
}
