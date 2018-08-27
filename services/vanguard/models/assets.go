package models

import (
	"fmt"

	"github.com/guregu/null"
)

type Assets struct {
	CharacterID   int32      `db:"characterID" json:"characterID"`
	CharacterName string     `db:"characterName" json:"characterName"`
	LocationFlag  string     `db:"locationFlag" json:"locationFlag"`
	LocationID    int64      `db:"locationID" json:"locationID"`
	TypeID        int64      `db:"typeID" json:"typeID"`
	ItemID        int64      `db:"itemID" json:"itemID"`
	TypeName      string     `db:"typeName" json:"typeName"`
	Quantity      int64      `db:"quantity" json:"quantity"`
	SubCount      int64      `db:"subCount" json:"subCount"`
	Buy           null.Float `db:"buy" json:"buy,omitempty"`
	Sell          null.Float `db:"sell" json:"sell,omitempty"`
	SubItems      []Assets   `db:"subItems" json:"subItems,omitempty"`
}

type AssetLocations struct {
	LocationFlag    string     `db:"locationFlag" json:"locationFlag"`
	LocationID      int64      `db:"locationID" json:"locationID"`
	LocationName    string     `db:"locationName" json:"locationName,omitempty"`
	SolarSystemName string     `db:"solarSystemName" json:"solarSystemName,omitempty"`
	Buy             null.Float `db:"buy" json:"buy,omitempty"`
	Sell            null.Float `db:"sell" json:"sell,omitempty"`
}

type AssetCharacters struct {
	CharacterID   int32      `db:"characterID" json:"characterID"`
	CharacterName string     `db:"characterName" json:"characterName"`
	Buy           null.Float `db:"buy" json:"buy,omitempty"`
	Sell          null.Float `db:"sell" json:"sell,omitempty"`
}

func GetAssetLocations(characterID int32, ownerHash string, filterCharacterID int32, marketable bool) ([]AssetLocations, error) {
	var filter string

	if filterCharacterID == 0 {
		filter = "IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = ''))"
	} else {
		// False AST, forced int64.
		filter = fmt.Sprintf("IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = '') AND tokenCharacterID=%d)", filterCharacterID)
	}

	if marketable {
		filter += " AND A.isSingleton = 0"
	}

	assetLocations := []AssetLocations{}
	if err := database.Select(&assetLocations, `
		SELECT A.locationID, stationName AS locationName, 
			SUM(P.sell  * IF(A.quantity, A.quantity, A.isSingleton)) AS sell
		FROM evedata.assets A
		JOIN evedata.jitaPrice P  ON A.typeID   = P.itemID
		JOIN staStations LOC ON LOC.stationID = A.locationID
		WHERE  A.characterID `+filter+`
		GROUP BY A.locationID
		ORDER BY sell DESC
	`, characterID, ownerHash); err != nil {
		return nil, err
	}
	return assetLocations, nil
}

func GetAssetCharacters(characterID int32, ownerHash string, marketable bool) ([]AssetCharacters, error) {
	filter := ""
	if marketable {
		filter = "A.isSingleton = 0 AND "
	}

	assetCharacters := []AssetCharacters{}
	if err := database.Select(&assetCharacters, `
		SELECT  A.characterID, characterName, 
			SUM(P.sell  * IF(A.quantity, A.quantity, A.isSingleton)) AS sell
		FROM evedata.assets A
		JOIN evedata.jitaPrice P  ON A.typeID   = P.itemID
		JOIN evedata.crestTokens C ON A.characterID = C.tokenCharacterID 
		WHERE `+filter+` A.characterID IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = ''))
		GROUP BY A.characterID
		ORDER BY sell DESC
	`, characterID, ownerHash); err != nil {
		return nil, err
	}
	return assetCharacters, nil
}

func GetAssets(characterID int32, ownerHash string, filterCharacterID int32, locationID int64) ([]Assets, error) {
	var filter string

	if filterCharacterID == 0 {
		filter = "IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = ''))"
	} else {
		// False AST, forced int64.
		filter = fmt.Sprintf("IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND (characterOwnerHash = ? OR characterOwnerHash = '') AND tokenCharacterID=%d)", filterCharacterID)
	}

	assets := []Assets{}
	if err := database.Select(&assets, `
		SELECT A.characterID, characterName, A.locationFlag, A.locationID, A.typeID, A.itemID,
			T.typeName, IF(A.quantity, A.quantity, A.isSingleton) AS quantity, 
			count(*) - 1 AS subCount,
			 -- Sum the layers of items
			 IFNULL(P.buy * IF(A.quantity, A.quantity, A.isSingleton),0) +
			 IFNULL(P1.buy * IF(L1.quantity, L1.quantity, L1.isSingleton),0)  +
			 IFNULL(P2.buy * IF(L2.quantity, L2.quantity, L2.isSingleton),0)  +
			 IFNULL(P3.buy * IF(L3.quantity, L3.quantity, L3.isSingleton),0)  AS buy, 
 
 			 IFNULL(P.sell * IF(A.quantity, A.quantity, A.isSingleton),0)  +
			 IFNULL(P1.sell * IF(L1.quantity, L1.quantity, L1.isSingleton),0)   +
			 IFNULL(P2.sell * IF(L2.quantity, L2.quantity, L2.isSingleton),0)   +
			 IFNULL(P3.sell * IF(L3.quantity, L3.quantity, L3.isSingleton),0)  AS sell

		FROM evedata.assets A
		JOIN evedata.crestTokens C on A.characterID = C.tokenCharacterID

		LEFT JOIN evedata.assets L1 ON A.itemID = L1.locationID
        LEFT JOIN evedata.assets L2 ON L1.itemID = L2.locationID
        LEFT JOIN evedata.assets L3 ON L1.itemID = L3.locationID
        
        -- Price everything in the 4 layers
		LEFT JOIN evedata.jitaPrice P ON A.typeID = P.itemID
		LEFT JOIN evedata.jitaPrice P1 ON L1.typeID = P1.itemID
		LEFT JOIN evedata.jitaPrice P2 ON L2.typeID = P2.itemID
		LEFT JOIN evedata.jitaPrice P3 ON L3.typeID = P3.itemID

		JOIN invTypes T ON A.typeID = T.typeID
		WHERE 
			A.characterID `+filter+`
			AND A.locationID = ?
		GROUP BY A.itemID
		ORDER BY sell DESC
	`, characterID, ownerHash, locationID); err != nil {
		return nil, err
	}

	count := 0
	errc := make(chan error)
	limit := make(chan bool, 10)
	for index := range assets {
		if assets[index].SubCount > 0 {
			count++
			go getSubAssets(assets[index].ItemID, &assets[index].SubItems, errc, limit)
		}
	}

	for i := 0; i < count; i++ {
		err := <-errc
		if err != nil {
			return nil, err // Something went wrong, break out.
		}
	}

	return assets, nil
}

func getSubAssets(itemID int64, assets *[]Assets, errc chan error, limit chan bool) {
	limit <- true
	defer func() { <-limit }()

	if err := database.Select(assets, `
		SELECT A.characterID, characterName, A.locationFlag, A.locationID, A.typeID, A.itemID,
			T.typeName, IF(A.quantity, A.quantity, A.isSingleton) AS quantity,
			count(*) - 1 AS subCount, 
			
			 IFNULL(P.buy * IF(A.quantity, A.quantity, A.isSingleton),0) +
			 IFNULL(P1.buy * IF(L1.quantity, L1.quantity, L1.isSingleton),0)  +
			 IFNULL(P2.buy * IF(L2.quantity, L2.quantity, L2.isSingleton),0)  +
			 IFNULL(P3.buy * IF(L3.quantity, L3.quantity, L3.isSingleton),0)  AS buy, 
 
 			 IFNULL(P.sell * IF(A.quantity, A.quantity, A.isSingleton),0)  +
			 IFNULL(P1.sell * IF(L1.quantity, L1.quantity, L1.isSingleton),0)   +
			 IFNULL(P2.sell * IF(L2.quantity, L2.quantity, L2.isSingleton),0)   +
			 IFNULL(P3.sell * IF(L3.quantity, L3.quantity, L3.isSingleton),0)  AS sell
			 
		FROM evedata.assets A
		JOIN evedata.crestTokens C on A.characterID = C.tokenCharacterID

		LEFT JOIN evedata.assets L1 ON A.itemID = L1.locationID
        LEFT JOIN evedata.assets L2 ON L1.itemID = L2.locationID
        LEFT JOIN evedata.assets L3 ON L1.itemID = L3.locationID
        
        -- Price everything in the 4 layers
		LEFT JOIN evedata.jitaPrice P ON A.typeID = P.itemID
		LEFT JOIN evedata.jitaPrice P1 ON L1.typeID = P1.itemID
		LEFT JOIN evedata.jitaPrice P2 ON L2.typeID = P2.itemID
		LEFT JOIN evedata.jitaPrice P3 ON L3.typeID = P3.itemID
		
		JOIN invTypes T ON A.typeID = T.typeID
		WHERE A.locationID = ?
		GROUP BY A.itemID
		ORDER BY sell DESC;
	`, itemID); err != nil {
		errc <- err
		return
	}

	count := 0
	a := *assets
	for index := range a {
		if a[index].SubCount > 0 {
			count++
			go getSubAssets(a[index].ItemID, &a[index].SubItems, errc, limit)
		}
	}

	for i := 0; i < count; i++ {
		err := <-errc
		if err != nil {
			errc <- err // Pass it on.
			return      // Something went wrong, break out.
		}
	}

	errc <- nil
}

type MarketableAssets struct {
	TypeID        int64      `db:"typeID" json:"typeID"`
	TypeName      string     `db:"typeName" json:"typeName"`
	Quantity      int64      `db:"quantity" json:"quantity"`
	Buy           null.Float `db:"buy" json:"buy,omitempty"`
	Sell          null.Float `db:"sell" json:"sell,omitempty"`
	StationPrice  null.Float `db:"stationPrice" json:"stationPrice,omitempty"`
	StationOrders null.Float `db:"stationOrders" json:"stationOrders,omitempty"`
	RegionPrice   null.Float `db:"regionPrice" json:"regionPrice,omitempty"`
	RegionOrders  null.Float `db:"regionOrders" json:"regionOrders,omitempty"`
}

func GetMarketableAssets(characterID int32, ownerHash string, tokenCharacterID int32, locationID int64) ([]MarketableAssets, error) {
	// Get the regionID
	var regionID int32
	if err := database.QueryRowx(`
		SELECT regionID FROM staStations WHERE stationID = ? LIMIT 1;`, locationID).Scan(&regionID); err != nil {
		return nil, err
	}
	assets := []MarketableAssets{}
	if err := database.Select(&assets, `
		SELECT A.typeID, typeName, A.quantity, buy, sell, 
			coalesce(stationPrice, 0) AS stationPrice, coalesce(regionPrice, 0) AS regionPrice,
			coalesce(regionOrders, 0) AS regionOrders, coalesce(stationOrders, 0) AS stationOrders 
			FROM evedata.assets A
			JOIN evedata.crestTokens C on A.characterID = C.tokenCharacterID	
			INNER JOIN evedata.jitaPrice J ON A.typeID = J.itemID
			INNER JOIN invTypes T ON A.typeID = T.typeID
			INNER JOIN staStations S ON S.stationID = A.locationID
			LEFT OUTER JOIN (SELECT typeID, count(*) AS regionOrders, min(price) AS regionPrice FROM evedata.market M FORCE INDEX(regionID_bid) WHERE regionID = ? AND bid = 0 GROUP BY typeID) MR ON MR.typeID = A.typeID
			LEFT OUTER JOIN (SELECT typeID, count(*) AS stationOrders, min(price) AS stationPrice FROM evedata.market M WHERE stationID = ? AND bid = 0 GROUP BY typeID) SR ON SR.typeID = A.typeID
			WHERE A.isSingleton = 0
				AND C.characterID = ?
				AND C.characterOwnerHash = ?
				AND A.characterID = ?
				AND A.locationID = ?
			ORDER BY typeName ASC
	`, regionID, locationID, characterID, ownerHash, tokenCharacterID, locationID); err != nil {
		return nil, err
	}

	return assets, nil
}
