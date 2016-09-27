package models

import (
	"evedata/null"
	"time"
)

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetActiveWarsByID(id int64) ([]CRESTRef, error) {
	w := []CRESTRef{}
	if err := database.Select(&w, `
			SELECT K.id, crestRef, type FROM
			(SELECT defenderID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM wars WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE (timeFinished = "0001-01-01 00:00:00" OR timeFinished IS NULL OR timeFinished >= UTC_TIMESTAMP()) AND timeStarted <= UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN crestID C ON C.id = K.id
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
			(SELECT defenderID AS id FROM wars WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM wars WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE timeStarted > timeDeclared AND timeStarted > UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN crestID C ON C.id = K.id
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
			(SELECT defenderID AS id FROM wars WHERE timeFinished < UTC_TIMESTAMP() AND aggressorID = ?
			UNION
			SELECT aggressorID AS id FROM wars WHERE timeFinished < UTC_TIMESTAMP() AND defenderID = ?
			UNION
			SELECT aggressorID  AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE timeFinished < UTC_TIMESTAMP() AND allyID = ?
			UNION
			SELECT allyID AS id FROM wars W INNER JOIN warAllies A on A.id = W.id WHERE timeFinished < UTC_TIMESTAMP() AND aggressorID = ?) AS K
			INNER JOIN crestID C ON C.id = K.id
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
	Kills         int64       `db:"kills" json:"kills"`
	Losses        int64       `db:"losses" json:"losses"`
}

// [BENCHMARK] 3.219 sec / 0.703 sec
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
	    IFNULL(kills,0) as kills,  
	    IFNULL(losses,0) as losses,
	    IF(AA.allianceID > 0, AA.name, AC.name) AS aggressorName,
	    IF(DA.allianceID > 0, DA.name, DC.name) AS defenderName
	        
		FROM wars W
		INNER JOIN crestID Ag ON Ag.id = aggressorID
	    INNER JOIN crestID Df ON Df.id = defenderID
	    LEFT OUTER JOIN alliances AA on AA.allianceID = aggressorID
		LEFT OUTER JOIN alliances DA on DA.allianceID = defenderID
		LEFT OUTER JOIN corporations AC on AC.corporationID = aggressorID
		LEFT OUTER JOIN corporations DC on DC.corporationID = defenderID
		LEFT OUTER JOIN 
	    ( -- Kills by the Aggressor
			SELECT 
				W.id, 
				count(*) AS kills
				FROM wars W
				INNER JOIN killmails K ON K.warID = W.id AND 
				(
					K.victimAllianceID != W.aggressorID AND 
					K.victimCorporationID != W.aggressorID
				)
				GROUP BY W.id
		) AS K ON W.id = K.id
		LEFT OUTER JOIN 
	    ( -- Kills by the Defenders
			SELECT 
				W.id, 
				count(*) AS losses
				FROM wars W
				INNER JOIN killmails L ON L.warID = W.id AND 
				(
					L.victimAllianceID = W.aggressorID OR 
					L.victimCorporationID = W.aggressorID
				)
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
	        
		FROM wars W
		INNER JOIN crestID Ag ON Ag.id = aggressorID
	    INNER JOIN crestID Df ON Df.id = defenderID
        LEFT OUTER JOIN warAllies A ON A.id = W.id
	    LEFT OUTER JOIN alliances AA on AA.allianceID = aggressorID
		LEFT OUTER JOIN alliances DA on DA.allianceID = defenderID
		LEFT OUTER JOIN corporations AC on AC.corporationID = aggressorID
		LEFT OUTER JOIN corporations DC on DC.corporationID = defenderID
		LEFT OUTER JOIN 
	    ( -- Kills by the Aggressor
			SELECT 
				W.id, 
				count(*) AS kills
				FROM wars W
				INNER JOIN killmails K ON K.warID = W.id AND 
				(
					K.victimAllianceID != W.aggressorID AND 
					K.victimCorporationID != W.aggressorID
				)
				GROUP BY W.id
		) AS K ON W.id = K.id
		LEFT OUTER JOIN 
	    ( -- Kills by the Defenders
			SELECT 
				W.id, 
				count(*) AS losses
				FROM wars W
				INNER JOIN killmails L ON L.warID = W.id AND 
				(
					L.victimAllianceID = W.aggressorID OR 
					L.victimCorporationID = W.aggressorID
				)
				GROUP BY W.id
		) AS L ON W.id = L.id
	    WHERE (aggressorID = ? OR defenderID = ? OR allyID = ?) AND
			(timeFinished > UTC_TIMESTAMP() OR
	        timeFinished = "0001-01-01 00:00:00")`, id, id, id); err != nil {
		return nil, err
	}
	return wars, nil
}
