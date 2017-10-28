package models

import (
	"fmt"

	"github.com/guregu/null"
)

type Assets struct {
	CharacterID   int64      `db:"characterID" json:"characterID"`
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
	CharacterID   int64      `db:"characterID" json:"characterID"`
	CharacterName string     `db:"characterName" json:"characterName"`
	Buy           null.Float `db:"buy" json:"buy,omitempty"`
	Sell          null.Float `db:"sell" json:"sell,omitempty"`
}

func GetAssetLocations(characterID int64, filterCharacterID int64) ([]AssetLocations, error) {
	var filter string

	if filterCharacterID == 0 {
		filter = "IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ?)"
	} else {
		// False AST, forced int64.
		filter = fmt.Sprintf("IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND tokenCharacterID=%d)", filterCharacterID)
	}

	assetLocations := []AssetLocations{}
	if err := database.Select(&assetLocations, `
		SELECT  A.locationID, stationName AS locationName, 
			SUM(P.sell  * IF(A.quantity, A.quantity, A.isSingleton)) AS sell
		FROM evedata.assets A
		JOIN evedata.jitaPrice P  ON A.typeID   = P.itemID
		JOIN staStations LOC ON LOC.stationID = A.locationID
		WHERE  A.characterID `+filter+`
		GROUP BY A.locationID
		ORDER BY sell DESC
	`, characterID); err != nil {
		return nil, err
	}
	return assetLocations, nil
}

func GetAssetCharacters(characterID int64) ([]AssetCharacters, error) {
	assetCharacters := []AssetCharacters{}
	if err := database.Select(&assetCharacters, `
		SELECT  A.characterID, characterName, 
			SUM(P.sell  * IF(A.quantity, A.quantity, A.isSingleton)) AS sell
		FROM evedata.assets A
		JOIN evedata.jitaPrice P  ON A.typeID   = P.itemID
		JOIN evedata.crestTokens C ON A.characterID = C.tokenCharacterID 
		WHERE A.characterID IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ?)
		GROUP BY A.characterID
		ORDER BY sell DESC
	`, characterID); err != nil {
		return nil, err
	}
	return assetCharacters, nil
}

// Obtain alliance information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetAssets(characterID int64, filterCharacterID int64, locationID int64) ([]Assets, error) {

	var filter string

	if filterCharacterID == 0 {
		filter = "IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ?)"
	} else {
		// False AST, forced int64.
		filter = fmt.Sprintf("IN (SELECT tokenCharacterID FROM evedata.crestTokens WHERE characterID = ? AND tokenCharacterID=%d)", filterCharacterID)
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
		WHERE A.locationType != "other"
			AND A.characterID `+filter+`
			AND A.locationID = ?
		GROUP BY A.itemID
		ORDER BY sell DESC
	`, characterID, locationID); err != nil {
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
