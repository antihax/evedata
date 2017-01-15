package models

import (
	"time"

	"github.com/guregu/null"
)

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetActiveWarsByID(id int64) ([]CRESTRef, error) {
	w := []CRESTRef{}
	if err := database.Select(&w, `
			SELECT K.id, crestRef, type FROM
			(SELECT defenderID AS id FROM evedata.wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM evedata.wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN evedata.crestID C ON C.id = K.id
		`, id, id, id, id); err != nil {
		return nil, err
	}
	return w, nil
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetPendingWarsByID(id int64) ([]CRESTRef, error) {
	w := []CRESTRef{}
	if err := database.Select(&w, `
			SELECT K.id, crestRef, type FROM
			(SELECT defenderID AS id FROM evedata.wars WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM evedata.wars WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN evedata.crestID C ON C.id = K.id
		`, id, id, id, id); err != nil {
		return nil, err
	}
	return w, nil
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetFinishedWarsByID(id int64) ([]CRESTRef, error) {
	w := []CRESTRef{}
	if err := database.Select(&w, `
			SELECT K.id, crestRef, type FROM
			(SELECT defenderID AS id FROM evedata.wars WHERE timeFinished < UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM evedata.wars WHERE timeFinished < UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE timeFinished < UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM evedata.wars W INNER JOIN evedata.warAllies A on A.id = W.id WHERE timeFinished < UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN evedata.crestID C ON C.id = K.id
		`, id, id, id, id); err != nil {
		return nil, err
	}
	return w, nil
}

type ActiveWarList struct {
	WarID         int64       `db:"warID" json:"warID"`
	TimeStarted   time.Time   `db:"timeStarted" json:"timeStarted"`
	TimeFinished  time.Time   `db:"timeFinished" json:"timeFinished"`
	OpenForAllies bool        `db:"openForAllies" json:"openForAllies"`
	AggressorID   int64       `db:"aggressorID" json:"aggressorID"`
	AggressorType null.String `db:"aggressorType" json:"aggressorType"`
	AggressorName null.String `db:"aggressorName" json:"aggressorName"`
	DefenderID    int64       `db:"defenderID" json:"defenderID"`
	DefenderType  null.String `db:"defenderType" json:"defenderType"`
	DefenderName  null.String `db:"defenderName" json:"defenderName"`
	Mutual        bool        `db:"mutual" json:"mutual"`
	WarKills      int64       `db:"warKills" json:"warKills"`
	WarLosses     int64       `db:"warLosses" json:"warLosses"`
	Efficiency    float64     `db:"efficiency" json:"efficiency"`
	Kills         int64       `db:"kills" json:"kills"`
	Losses        int64       `db:"losses" json:"losses"`
}

// [BENCHMARK] 1.469 sec / 0.094 sec
func GetActiveWarList() ([]ActiveWarList, error) {
	wars := []ActiveWarList{}
	if err := database.Select(&wars, `
	SELECT 
		W.id AS warID, 
	    timeStarted, 
	    timeFinished, 
	    openForAllies, 
	    aggressorID, 
	    Ag.Type AS aggressorType, 
	    defenderID, 
	    Df.type AS defenderType, 
	    mutual, 
	    IFNULL(K.kills,0) as warKills,  
	    IFNULL(L.losses,0) as warLosses,
	    IF(AA.allianceID > 0, AA.name, AC.name) AS aggressorName,
	    IF(DA.allianceID > 0, DA.name, DC.name) AS defenderName,
		IFNULL(S.efficiency,1) AS efficiency,
        IFNULL(S.kills,0) AS kills,
        IFNULL(S.losses,0)  AS losses
		FROM evedata.wars W
		INNER JOIN evedata.crestID Ag ON Ag.id = aggressorID
	    INNER JOIN evedata.crestID Df ON Df.id = defenderID
	    LEFT OUTER JOIN evedata.alliances AA on AA.allianceID = aggressorID
		LEFT OUTER JOIN evedata.alliances DA on DA.allianceID = defenderID
		LEFT OUTER JOIN evedata.corporations AC on AC.corporationID = aggressorID
		LEFT OUTER JOIN evedata.corporations DC on DC.corporationID = defenderID
        LEFT OUTER JOIN evedata.entityKillStats S ON S.id = aggressorID
		LEFT OUTER JOIN 
	    ( -- Kills by the Aggressor
			SELECT 
				W.id, 
				count(*) AS kills
				FROM evedata.wars W
				INNER JOIN evedata.killmails K ON K.warID = W.id AND 
				(
					K.victimAllianceID != W.aggressorID AND 
					K.victimCorporationID != W.aggressorID
				)
				WHERE killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 31 DAY)
				GROUP BY W.id
		) AS K ON W.id = K.id
		LEFT OUTER JOIN 
	    ( -- Kills by the Defenders
			SELECT 
				W.id, 
				count(*) AS losses
				FROM evedata.wars W
				INNER JOIN evedata.killmails L ON L.warID = W.id AND 
				(
					L.victimAllianceID = W.aggressorID OR 
					L.victimCorporationID = W.aggressorID
				)
				WHERE killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 31 DAY)
				GROUP BY W.id
		) AS L ON W.id = L.id
	    WHERE mutual = 0 AND
			(timeFinished > UTC_TIMESTAMP() OR
	        timeFinished = "0001-01-01 00:00:00")`); err != nil {
		return nil, err
	}
	return wars, nil
}

func GetWarsForEntityByID(id int64) ([]ActiveWarList, error) {
	wars := []ActiveWarList{}
	if err := database.Select(&wars, `
	SELECT 
		W.id AS warID, 
	    timeStarted, 
	    timeFinished, 
	    openForAllies, 
	    aggressorID, 
	    Ag.Type AS aggressorType, 
	    defenderID, 
	    Df.type AS defenderType, 
	    mutual, 
	    IFNULL(kills,0) as kills,  
	    IFNULL(losses,0) as losses,
	    IF(AA.allianceID > 0, AA.name, AC.name) AS aggressorName,
	    IF(DA.allianceID > 0, DA.name, DC.name) AS defenderName
	        
		FROM evedata.wars W
		INNER JOIN evedata.crestID Ag ON Ag.id = aggressorID
	    INNER JOIN evedata.crestID Df ON Df.id = defenderID
        LEFT OUTER JOIN evedata.warAllies A ON A.id = W.id
	    LEFT OUTER JOIN evedata.alliances AA on AA.allianceID = aggressorID
		LEFT OUTER JOIN evedata.alliances DA on DA.allianceID = defenderID
		LEFT OUTER JOIN evedata.corporations AC on AC.corporationID = aggressorID
		LEFT OUTER JOIN evedata.corporations DC on DC.corporationID = defenderID
		LEFT OUTER JOIN 
	    ( -- Kills by the Aggressor
			SELECT 
				W.id, 
				count(*) AS kills
				FROM evedata.wars W
				INNER JOIN evedata.killmails K ON K.warID = W.id AND 
				(
					K.victimAllianceID != W.aggressorID AND 
					K.victimCorporationID != W.aggressorID
				)
				WHERE killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 31 DAY)
				GROUP BY W.id
		) AS K ON W.id = K.id
		LEFT OUTER JOIN 
	    ( -- Kills by the Defenders
			SELECT 
				W.id, 
				count(*) AS losses
				FROM evedata.wars W
				INNER JOIN evedata.killmails L ON L.warID = W.id AND 
				(
					L.victimAllianceID = W.aggressorID OR 
					L.victimCorporationID = W.aggressorID
				)
				WHERE killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 31 DAY)
				GROUP BY W.id
		) AS L ON W.id = L.id
	    WHERE (aggressorID = ? OR defenderID = ? OR allyID = ?) AND
			(timeFinished > UTC_TIMESTAMP() OR
	        timeFinished = "0001-01-01 00:00:00")`, id, id, id); err != nil {
		return nil, err
	}
	return wars, nil
}

type KnownAllies struct {
	Number int64  `db:"number" json:"number"`
	AllyID int64  `db:"allyID" json:"allyID"`
	Name   string `db:"name" json:"name"`
	Type   string `db:"type" json:"type"`
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetKnownAlliesByID(id int64) ([]KnownAllies, error) {
	w := []KnownAllies{}
	if err := database.Select(&w, `
			SELECT 
				COUNT(DISTINCT W.id) AS number, 
			    allyID, 
			    CREST.type,
				IFNULL(DA.name, DC.name) AS name
			FROM evedata.wars W
				INNER JOIN evedata.warAllies A ON W.id = A.id
				INNER JOIN evedata.crestID CREST ON CREST.id = A.allyID
				LEFT OUTER JOIN evedata.alliances DA on DA.allianceID = A.allyID
				LEFT OUTER JOIN evedata.corporations DC on DC.corporationID = A.allyID
				WHERE defenderID = ? AND W.timeStarted > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 12 MONTH)
				GROUP BY allyID
		`, id); err != nil {
		return nil, err
	}
	return w, nil
}
