package models

type MarketItemList struct {
	TypeID     int64  `db:"typeID"`
	TypeName   string `db:"typeName"`
	Categories string `db:"Categories"`
	Count      int64
}

func SearchMarketNames(query string) ([]MarketItemList, error) {
	mIL := []MarketItemList{}

	// [BENCHMARK] 0.078 sec / 0.000 sec
	err := database.Select(&mIL, `SELECT  T.typeID, typeName, CONCAT_WS(',', G5.marketGroupName, G4.marketGroupName, G3.marketGroupName, G2.marketGroupName, G.marketGroupName) AS Categories, count(*) AS count
           FROM invTypes T 
           LEFT JOIN invMarketGroups G on T.marketGroupID = G.marketGroupID
           LEFT JOIN invMarketGroups G2 on G.parentGroupID = G2.marketGroupID
           LEFT JOIN invMarketGroups G3 on G2.parentGroupID = G3.marketGroupID
           LEFT JOIN invMarketGroups G4 on G3.parentGroupID = G4.marketGroupID
           LEFT JOIN invMarketGroups G5 on G4.parentGroupID = G5.marketGroupID

           WHERE published=1 AND T.marketGroupID IS NOT NULL AND typeName LIKE ?
           GROUP BY T.typeID
           ORDER BY typeName
           LIMIT 100`, "%"+query+"%")
	if err != nil {
		return nil, err
	}

	return mIL, nil
}
