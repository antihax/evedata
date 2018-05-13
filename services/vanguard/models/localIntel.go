package models

import "strings"

type LocalIntelData struct {
	Number         int64   `db:"number" json:"number"`
	ID             int64   `db:"id" json:"id"`
	EntityName     string  `db:"entityName" json:"entityName"`
	FactionName    string  `db:"factionName" json:"factionName"`
	Type           string  `db:"type" json:"type"`
	MemberCount    int64   `db:"memberCount" json:"memberCount"`
	WarAggressor   int64   `db:"warAggressor" json:"warAggressor"`
	WarDefender    int64   `db:"warDefender" json:"warDefender"`
	Kills          int64   `db:"kills" json:"kills"`
	Losses         int64   `db:"losses" json:"losses"`
	CapKills       int64   `db:"capKills" json:"capKills"`
	Efficiency     float64 `db:"efficiency" json:"efficiency"`
	CapProbability float64 `db:"capProbability" json:"capProbability"`
}

// [BENCHMARK] 1.469 sec / 0.094 sec
// FALSE Positive AST, concatenation is for replace tokens, actual values are fed
// through vargs.
func GetLocalIntel(names []interface{}) ([]LocalIntelData, error) {
	wars := []LocalIntelData{}
	if err := database.Select(&wars, `
		SELECT 	number,
				entityName,
				SUB1.id,
				type,
				memberCount,
				IFNULL(factionName, "") AS factionName,
				COUNT(DISTINCT Agg.ID) AS warAggressor,
				COUNT(DISTINCT Def.ID) AS warDefender,
				COALESCE(kills,0) AS kills,
				COALESCE(losses,0) AS losses,
				COALESCE(capKills,0) AS capKills,
				IF(kills, capKills / kills, 0) AS capProbability,
				IF(losses+kills, (kills/(kills+losses)), 1.0000) AS efficiency
		FROM
			(SELECT
				COUNT(DISTINCT Ch.characterID) AS number,
				IF(A.allianceID, A.name, Co.name) AS entityName,
				CREST.id,
				itemName AS factionName,
				CREST.type,
				IF(A.allianceID, A.memberCount, Co.memberCount) AS memberCount,
				sum(kills) AS kills,
				sum(losses) AS losses,
				sum(capKills) AS capKills 
				FROM evedata.characters Ch
				LEFT OUTER JOIN evedata.alliances A ON Ch.allianceID = A.allianceID
				LEFT OUTER JOIN evedata.corporations Co ON Ch.corporationID = Co.corporationID
				LEFT OUTER JOIN invNames Fa ON Fa.itemID = Co.factionID
				LEFT OUTER JOIN evedata.entityKillStats S ON S.id = Ch.characterID
				INNER JOIN evedata.entities CREST ON CREST.id = IF(A.allianceID, A.allianceID, Co.corporationID)
				WHERE Ch.name IN (?`+strings.Repeat(",?", len(names)-1)+`)
				GROUP BY CREST.id) SUB1
		LEFT OUTER JOIN evedata.wars Agg ON Agg.aggressorID = SUB1.id AND (Agg.timeFinished = "0001-01-01 00:00:00" OR Agg.timeFinished IS NULL OR Agg.timeFinished >= UTC_TIMESTAMP()) AND Agg.timeStarted <= UTC_TIMESTAMP()
		LEFT OUTER JOIN evedata.wars Def ON Def.defenderID = SUB1.id AND (Def.timeFinished = "0001-01-01 00:00:00" OR Def.timeFinished IS NULL OR Def.timeFinished >= UTC_TIMESTAMP()) AND Def.timeStarted <= UTC_TIMESTAMP()
		GROUP BY SUB1.id`, names...); err != nil {
		return nil, err
	}
	return wars, nil
}
