package models

import "github.com/guregu/null"

type Assets struct {
	CharacterID  int64      `db:"characterID" json:"characterID"`
	LocationFlag string     `db:"locationFlag" json:"locationFlag"`
	LocationID   int64      `db:"locationID" json:"locationID"`
	TypeID       int64      `db:"typeID" json:"typeID"`
	ItemID       int64      `db:"itemID" json:"itemID"`
	TypeName     string     `db:"typeName" json:"typeName"`
	Quantity     int64      `db:"quantity" json:"quantity"`
	SubCount     int64      `db:"subCount" json:"subCount"`
	Buy          null.Float `db:"buy" json:"buy,omitempty"`
	Sell         null.Float `db:"sell" json:"sell,omitempty"`
	SubItems     []Assets   `db:"subItems" json:"subItems,omitempty"`
}

type AssetLocations struct {
	LocationFlag    string     `db:"locationFlag" json:"locationFlag"`
	LocationID      int64      `db:"locationID" json:"locationID"`
	LocationName    string     `db:"locationName" json:"locationName,omitempty"`
	SolarSystemName string     `db:"solarSystemName" json:"solarSystemName,omitempty"`
	Buy             null.Float `db:"buy" json:"buy,omitempty"`
	Sell            null.Float `db:"sell" json:"sell,omitempty"`
}

func GetAssetLocations(characterID int64, filterCharacterID int64) ([]AssetLocations, error) {
	assetLocations := []AssetLocations{}
	if err := database.Select(&assetLocations, `
		SELECT  A.locationFlag, A.locationID, stationName as locationName,
			 -- Sum the layers of items
			 SUM(IFNULL(P.buy,0))  + SUM(IFNULL(P1.buy,0))  + SUM(IFNULL(P2.buy,0))  + SUM(IFNULL(P3.buy,0))  AS buy, 
             SUM(IFNULL(P.sell,0)) + SUM(IFNULL(P1.sell,0)) + SUM(IFNULL(P2.sell,0)) + SUM(IFNULL(P3.sell,0)) AS sell

        -- Work through 4 layers
		FROM evedata.assets A
		LEFT JOIN evedata.assets L1 ON A.itemID = L1.locationID
        LEFT JOIN evedata.assets L2 ON L1.itemID = L2.locationID
        LEFT JOIN evedata.assets L3 ON L1.itemID = L3.locationID
        
        -- Price everything in the 4 layers
		LEFT JOIN evedata.jitaPrice P ON A.typeID = P.itemID
		LEFT JOIN evedata.jitaPrice P1 ON L1.typeID = P1.itemID
		LEFT JOIN evedata.jitaPrice P2 ON L2.typeID = P2.itemID
		LEFT JOIN evedata.jitaPrice P3 ON L3.typeID = P3.itemID
 
		JOIN invTypes T ON A.typeID = T.typeID
        LEFT JOIN staStations LOC ON LOC.stationID = A.locationID
		WHERE A.locationType != "other"
			AND A.characterID IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ?)
		GROUP BY A.locationID
	`, characterID); err != nil {
		return nil, err
	}
	return assetLocations, nil
}

// Obtain alliance information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetAssets(characterID int64) ([]Assets, error) {
	assets := []Assets{}
	if err := database.Select(&assets, `
		SELECT A.characterID, A.locationFlag, A.locationID, A.typeID, A.itemID,
			T.typeName, IF(A.quantity, A.quantity, A.isSingleton) AS quantity, 
			count(*) - 1 AS subCount, buy, sell
			
		FROM evedata.assets A
		LEFT JOIN evedata.assets L ON A.itemID = L.locationID
		LEFT JOIN evedata.jitaPrice P ON A.typeID = P.itemID
		JOIN invTypes T ON A.typeID = T.typeID
		WHERE A.locationType != "other"
			AND A.characterID IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ?)
		GROUP BY A.locationID, A.itemID;
	`, characterID); err != nil {
		return nil, err
	}

	count := 0
	errc := make(chan error)
	limit := make(chan bool, 25)
	for index, _ := range assets {
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
		SELECT A.characterID, A.locationFlag, A.locationID, A.typeID, A.itemID,
			T.typeName, IF(A.quantity, A.quantity, A.isSingleton) AS quantity, 
			count(*) - 1 AS subCount, buy, sell
			
		FROM evedata.assets A
		LEFT JOIN evedata.assets L ON A.itemID = L.locationID
		LEFT JOIN evedata.jitaPrice P ON A.typeID = P.itemID
		JOIN invTypes T ON A.typeID = T.typeID
		WHERE A.locationID = ?
		GROUP BY A.locationID, A.itemID
		ORDER BY A.locationID;
	`, itemID); err != nil {
		errc <- err
		return
	}

	count := 0
	a := *assets
	for index, _ := range a {
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
