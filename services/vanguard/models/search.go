package models

type MarketItemList struct {
	TypeID     int64  `db:"typeID"`
	TypeName   string `db:"typeName"`
	Categories string `db:"Categories"`
	Count      int64
}

func SearchMarketNames(query string) ([]MarketItemList, error) {
	list := []MarketItemList{}

	// [BENCHMARK] 0.078 sec / 0.000 sec
	err := database.Select(&list, `SELECT  T.typeID, typeName, CONCAT_WS(',', G5.marketGroupName, G4.marketGroupName, G3.marketGroupName, G2.marketGroupName, G.marketGroupName) AS Categories, count(*) AS count
           FROM invTypes T 
           LEFT JOIN invMarketGroups G on T.marketGroupID = G.marketGroupID
           LEFT JOIN invMarketGroups G2 on G.parentGroupID = G2.marketGroupID
           LEFT JOIN invMarketGroups G3 on G2.parentGroupID = G3.marketGroupID
           LEFT JOIN invMarketGroups G4 on G3.parentGroupID = G4.marketGroupID
           LEFT JOIN invMarketGroups G5 on G4.parentGroupID = G5.marketGroupID

           WHERE published=1 AND T.marketGroupID IS NOT NULL AND typeName LIKE ?
           GROUP BY T.typeID
           ORDER BY typeName
           LIMIT 100`, query+"%")
	if err != nil {
		return nil, err
	}

	return list, nil
}

type NamesItemList struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Type string `db:"type" json:"type"`
}

func SearchNames(query string) ([]NamesItemList, error) {
	list := []NamesItemList{}
	query = query + "%"

	// [BENCHMARK] 0.000 sec / 0.000 sec
	err := database.Select(&list, `
		SELECT typeName AS name, typeID AS id, "Item" AS type 
			FROM invTypes WHERE typeName LIKE ?
			UNION
			SELECT name, characterID AS id, "Character" AS type 
			FROM evedata.characters WHERE name LIKE ?
			UNION
			SELECT name, corporationID AS id, "Corporation" AS type 
			FROM evedata.corporations WHERE name LIKE ?
			UNION
			SELECT name, allianceID AS id, "Alliance" AS type 
			FROM evedata.alliances WHERE name LIKE ?
			ORDER BY name ASC;
		`, query, query, query, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func SearchEntities(query string) ([]NamesItemList, error) {
	list := []NamesItemList{}
	query = query + "%"

	// [BENCHMARK] 0.000 sec / 0.000 sec
	err := database.Select(&list, `
		SELECT name, id, type FROM (
			SELECT name, corporationID AS id, "Corporation" AS type 
			FROM evedata.corporations WHERE name LIKE ?
			UNION
			SELECT name, allianceID AS id, "Alliance" AS type 
			FROM evedata.alliances WHERE name LIKE ?) A
			ORDER BY name ASC;
		`, query, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}
