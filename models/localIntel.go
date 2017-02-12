package models

import "strings"

type LocalIntelData struct {
	Number       int64   `db:"number" json:"number"`
	ID           int64   `db:"id" json:"id"`
	EntityName   string  `db:"entityName" json:"entityName"`
	Type         string  `db:"type" json:"type"`
	MemberCount  int64   `db:"memberCount" json:"memberCount"`
	WarAggressor int64   `db:"warAggressor" json:"warAggressor"`
	WarDefender  int64   `db:"warDefender" json:"warDefender"`
	Kills        int64   `db:"kills" json:"kills"`
	Losses       int64   `db:"losses" json:"losses"`
	Efficiency   float64 `db:"efficiency" json:"efficiency"`
}

// [BENCHMARK] 1.469 sec / 0.094 sec
// FALSE Positive AST, concatenation is for replace tokens, actual values are fed
// through vargs.
func GetLocalIntel(names []interface{}) ([]LocalIntelData, error) {
	wars := []LocalIntelData{}
	if err := database.Select(&wars, `
			   SELECT number,
				   entityName,
				   SUB1.id,
				   type,
				   memberCount,
			       COUNT(DISTINCT Agg.ID) AS warAggressor,
			       COUNT(DISTINCT Def.ID) AS warDefender,
		           kills,
		           losses,
		           efficiency
			   FROM
			   	(SELECT
			   		COUNT(DISTINCT Ch.characterID) AS number,
			   		IF(A.allianceID, A.name, Co.name) AS entityName,
			   		CREST.id,
			   		CREST.type,
			           IF(A.allianceID, A.memberCount, Co.memberCount) AS memberCount
			   	FROM evedata.characters Ch
			   	LEFT OUTER JOIN evedata.alliances A ON Ch.allianceID = A.allianceID

			   	LEFT OUTER JOIN evedata.corporations Co ON Ch.corporationID = Co.corporationID
			   	INNER JOIN evedata.crestID CREST ON CREST.id = IF(A.allianceID, A.allianceID, Co.corporationID)

			   	WHERE Ch.name IN (?`+strings.Repeat(",?", len(names)-1)+`)
			   	GROUP BY CREST.id) SUB1
				LEFT OUTER JOIN evedata.entityKillStats S ON S.id = SUB1.id
			   LEFT OUTER JOIN evedata.wars Agg ON Agg.aggressorID = SUB1.id AND (Agg.timeFinished = "0001-01-01 00:00:00" OR Agg.timeFinished IS NULL OR Agg.timeFinished >= UTC_TIMESTAMP()) AND Agg.timeStarted <= UTC_TIMESTAMP()
			   LEFT OUTER JOIN evedata.wars Def ON Def.defenderID = SUB1.id AND (Def.timeFinished = "0001-01-01 00:00:00" OR Def.timeFinished IS NULL OR Def.timeFinished >= UTC_TIMESTAMP()) AND Def.timeStarted <= UTC_TIMESTAMP()
			   WHERE kills IS NOT NULL AND losses IS NOT NULL AND efficiency IS NOT NULL
			   GROUP BY SUB1.id`, names...); err != nil {
		return nil, err
	}
	return wars, nil
}
