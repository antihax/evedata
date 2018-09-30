package models

import (
	"database/sql"
	"fmt"
)

type MarketOrderCharacters struct {
	CharacterID   int32  `db:"characterID" json:"characterID"`
	CharacterName string `db:"characterName" json:"characterName"`
}

func GetOrderCharacters(characterID int32, ownerHash string) ([]MarketOrderCharacters, error) {
	c := []MarketOrderCharacters{}
	if err := database.Select(&c, `
		SELECT  O.characterID, characterName
		FROM evedata.orders O
		JOIN evedata.crestTokens C ON O.characterID = C.tokenCharacterID 
		WHERE O.characterID IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = ''))
		GROUP BY O.characterID
		ORDER BY characterName DESC
	`, characterID, ownerHash); err != nil {
		return nil, err
	}
	return c, nil
}

type MarketOrders struct {
	TypeID        int64           `db:"typeID" json:"typeID"`
	TypeName      string          `db:"typeName" json:"typeName"`
	CharacterID   int64           `db:"characterID" json:"characterID"`
	CharacterName string          `db:"characterName" json:"characterName"`
	StationName   string          `db:"stationName" json:"stationName"`
	StationID     int64           `db:"stationID" json:"stationID"`
	OrderID       int64           `db:"orderID" json:"orderID"`
	RegionName    string          `db:"regionName" json:"regionName"`
	LocationID    int64           `db:"locationID" json:"locationID"`
	IsBuyOrder    int64           `db:"isBuyOrder" json:"isBuyOrder"`
	VolumeRemain  int64           `db:"volumeRemain" json:"volumeRemain"`
	VolumeTotal   int64           `db:"volumeTotal" json:"volumeTotal"`
	Price         sql.NullFloat64 `db:"price" json:"price"`
	CurrentPrice  sql.NullFloat64 `db:"currentPrice" json:"currentPrice"`
	RegionPrice   sql.NullFloat64 `db:"regionPrice" json:"regionPrice"`
	NumOrders     int64           `db:"numOrders" json:"numOrders"`
}

// GetOrders for a character
func GetOrders(characterID int32, ownerHash string, filterCharacterID int32) ([]MarketOrders, error) {
	var filter string

	if filterCharacterID == 0 {
		filter = "IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = ''))"
	} else {
		// False AST, forced int64.
		filter = fmt.Sprintf("IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = '') AND tokenCharacterID=%d)", filterCharacterID)
	}

	c := []MarketOrders{}
	if err := database.Select(&c, `
	SELECT 	O.characterID, C.characterName,S.stationID,  S.stationName,
			O.orderID, R.regionName, O.locationID, O.isBuyOrder,
			O.typeID, T.typeName,
			O.volumeRemain, O.volumeTotal, O.price, 
			IF(O.isBuyOrder, max(M.price), min(M.price)) AS currentPrice,
			IF(O.isBuyOrder, max(MR.price), min(MR.price)) AS regionPrice,
			count(*) AS numOrders
	FROM evedata.orders O
	INNER JOIN eve.invTypes T ON T.typeID = O.typeID
	INNER JOIN evedata.crestTokens C ON C.tokenCharacterID = O.characterID
	INNER JOIN evedata.structures S ON S.stationID = O.locationID
	INNER JOIN eve.mapSolarSystems MS ON MS.solarSystemID = S.solarSystemID
	INNER JOIN eve.mapRegions R ON R.regionID = MS.regionID
	LEFT OUTER JOIN evedata.market M FORCE INDEX(ix_marketorders) ON 
		M.typeID = O.typeID AND
		M.bid = O.isBuyOrder AND
		M.stationID = O.locationID AND
		M.orderID != O.orderID
	LEFT OUTER JOIN evedata.market MR FORCE INDEX(ix_marketordersregion) ON 
		MR.typeID = O.typeID AND
		MR.bid = O.isBuyOrder AND
		MR.regionID = R.regionID AND 
		MR.orderID != O.orderID
	WHERE 
		O.characterID `+filter+`
	GROUP BY O.orderID

	`, characterID, ownerHash); err != nil {
		return nil, err
	}

	return c, nil
}
