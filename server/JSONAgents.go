package evedata

import (
	"encoding/json"
	"errors"
	"evedata-revel/null"
	"evedata/models"
	"net/http"
	"regexp"
	"strconv"
)

func init() {
	AddRoute(Route{"agents", "GET", "/U/agents", FindAgents})
}

/******************************************************************************
 * marketRegions JSON query
 *****************************************************************************/
type agents struct {
	AgentID       int64       `db:"agentID"       json:"agentID"`
	Level         int64       `db:"level"         json:"level"`
	Agent         string      `db:"agent"         json:"agent"`
	Corp          string      `db:"corp"          json:"corp"`
	Division      string      `db:"division"      json:"division"`
	Faction       string      `db:"faction"       json:"faction"`
	Security      string      `db:"security"      json:"security"`
	SolarSystemID int64       `db:"solarSystemID" json:"solarSystemID"`
	StationID     int64       `db:"stationID"     json:"stationID"`
	Station       string      `db:"station"       json:"station"`
	Required      null.String `db:"RequiredStanding"`
	AgentStd      null.String `db:"AgentStanding"`
	CorpStnd      null.String `db:"CorpStanding"`
	FactionStd    null.String `db:"FactionStanding"`
	Jumps         string      `db:"J"             json:"jumps"`
}

// ARows bridge for old version
type ARows struct {
	Rows *[]agents `json:"rows"`
}

// FindAgents generate a list of agents based on user input
func FindAgents(c *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	aG := []agents{}

	var (
		err         error
		mRows       ARows
		systemID    int
		characterID int
		level       int
		jumps       int
		sec         string
		highsec     bool
		locator     bool
		locate      string
		division    int
		divisionS   string
	)

	// Validate user input

	systemID, err = strconv.Atoi(r.FormValue("systemID"))
	if err != nil {
		return 500, errors.New("Invalid systemID")
	}

	characterID, err = strconv.Atoi(r.FormValue("characterID"))

	level, err = strconv.Atoi(r.FormValue("level"))
	if err != nil {
		return 500, errors.New("Invalid level")
	}

	jumps, err = strconv.Atoi(r.FormValue("jumps"))
	if err != nil {
		return 500, errors.New("Invalid jumps")
	}

	match, err := regexp.MatchString("([0-9].[0-9])", r.FormValue("sec"))
	if err != nil {
		return 500, errors.New("Invalid sec")
	}
	if match == true {
		sec = r.FormValue("sec")
	}

	highsec = BooleanizeFormValue(r.FormValue("highsec"))
	locator = BooleanizeFormValue(r.FormValue("locator"))

	division, err = strconv.Atoi(r.FormValue("division"))

	// Build custom strings for query filters

	if locator == true {
		locate = "		         Cfg.k = \"agent.LocateCharacterService.enabled\" AND "
	} else {
		locate = ""
	}

	if division > 0 {
		divisionS = "		         A.divisionID = " + strconv.Itoa(division) + "  AND "
	} else {
		divisionS = ""
	}

	user := models.GetUser(r)

	// SECURED: requires user to be logged in, once we know that we add
	// restriction to the query to ensure only linked cid -> uid can obtain data
	if characterID != 0 && user != nil {
		var sqlQuery string
		sqlQuery = `
		SELECT
		         A.agentID,
		         A.level,
		         E.itemName AS agent,
		         Ce.itemName AS corp,
		         CD.divisionName as division,
		         Fe.itemName AS faction,
		         ROUND(Sys.security,1) as security,
		         Sta.stationID,
		         Sta.stationName as station,
		         Sta.solarSystemID,
		         ((A.level - 1) * 2 + -20 / 20) AS RequiredStanding,
		         AgS.standing + (10-AgS.standing)*(0.04*IF(AgS.standing>0,Con.level,Dip.level)) AS AgentStanding,
		         CoS.standing + (10-CoS.standing)*(0.04*IF(CoS.standing>0,Con.level,Dip.level)) AS CorpStanding,
		         FaS.standing + (10-FaS.standing)*(0.04*IF(FaS.standing>0,Con.level,Dip.level)) AS FactionStanding,`

		if highsec == true {
			sqlQuery += `secureJumps AS J`
		} else {
			sqlQuery += `jumps AS J`
		}

		sqlQuery += `		
			FROM agtAgents AS A
		         INNER JOIN agtConfig AS Cfg ON A.agentID = Cfg.agentID
		         INNER JOIN eveNames AS E ON A.agentID = E.itemID
		         INNER JOIN staStations AS Sta ON Sta.stationID = A.locationID
		         INNER JOIN mapSolarSystems AS Sys ON Sta.solarSystemID = Sys.solarSystemID
		         INNER JOIN crpNPCCorporations AS C ON A.corporationID = C.corporationID
		         INNER JOIN crpNPCDivisions AS CD ON A.divisionID = CD.divisionID
		         INNER JOIN eveNames AS Ce ON A.corporationID = Ce.itemID
		         INNER JOIN eveNames AS Fe ON C.factionID = Fe.itemID
		         INNER JOIN jumps AS Jm ON Jm.toSolarSystemID = Sys.solarSystemID
		         INNER JOIN charSkills AS Dip ON Dip.cid = ? AND Dip.typeID = 3357
		         INNER JOIN charSkills AS Con ON Con.cid = ? AND Con.typeID = 3359
		         LEFT OUTER JOIN standings AS AgS ON (A.agentID = AgS.itemID) AND AgS.cid = ?
		         LEFT OUTER JOIN standings AS CoS ON (A.corporationID = CoS.itemID) AND CoS.cid = ?
		         LEFT OUTER JOIN standings AS FaS ON (C.factionID = FaS.itemID) AND FaS.cid = ?
		         INNER JOIN characters AS Cx ON Cx.characterID = ? AND Cx.uid = ?
		WHERE`

		sqlQuery += locate
		sqlQuery += divisionS

		sqlQuery += `         A.level >= ? AND
		         Sys.security <= ? AND
		         Jm.fromSolarSystemID = ?
		GROUP BY
		         A.agentID
		HAVING
		         J > 0 AND
		         J <= ? AND (
		            CorpStanding > RequiredStanding OR
		            FactionStanding > RequiredStanding OR
		            AgentStanding > RequiredStanding
         )
		ORDER BY J;`

		err = c.Db.Select(&aG, sqlQuery,
			characterID, characterID, characterID, characterID, characterID, characterID,
			user.UID, level, sec, systemID, jumps)

	} else {
		var sqlQuery string
		sqlQuery = `
		SELECT
		         A.agentID,
		         A.level,
		         E.itemName AS agent,
		         Ce.itemName AS corp,
		         CD.divisionName as division,
		         Fe.itemName AS faction,
		         ROUND(Sys.security,1) as security,
		         Sta.stationID,
		         Sta.stationName as station,
		         Sta.solarSystemID,
		         ((A.level - 1) * 2 + -20 / 20) AS RequiredStanding,`
		if highsec == true {
			sqlQuery += `secureJumps AS J`
		} else {
			sqlQuery += `jumps AS J`
		}

		sqlQuery += `		
				 FROM agtAgents AS A
		         INNER JOIN agtConfig AS Cfg ON A.agentID = Cfg.agentID
		         INNER JOIN eveNames AS E ON A.agentID = E.itemID
		         INNER JOIN staStations AS Sta ON Sta.stationID = A.locationID
		         INNER JOIN mapSolarSystems AS Sys ON Sta.solarSystemID = Sys.solarSystemID
		         INNER JOIN crpNPCCorporations AS C ON A.corporationID = C.corporationID
		         INNER JOIN crpNPCDivisions AS CD ON A.divisionID = CD.divisionID
		         INNER JOIN eveNames AS Ce ON A.corporationID = Ce.itemID
		         INNER JOIN eveNames AS Fe ON C.factionID = Fe.itemID
		         INNER JOIN jumps AS Jm ON Jm.toSolarSystemID = Sys.solarSystemID
		WHERE`

		sqlQuery += locate
		sqlQuery += divisionS

		sqlQuery += `         A.level >= ? AND
		         Sys.security <= ? AND
		         Jm.fromSolarSystemID = ?
		GROUP BY
		         A.agentID
		HAVING
		         J > 0 AND
		         J <= ?
		ORDER BY J         ;`

		err = c.Db.Select(&aG, sqlQuery, level, sec, systemID, jumps)
	}
	mRows.Rows = &aG

	if err != nil {
		return 500, err
	}

	// Skip the root node and JSONify.
	encoder := json.NewEncoder(w)
	encoder.Encode(mRows)

	return 200, nil
}
