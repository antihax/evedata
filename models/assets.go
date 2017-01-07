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

// Obtain alliance information by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetAssets(characterID int64) ([]Assets, error) {
	assets := []Assets{}
	if err := database.Select(&assets, `
		SELECT A.characterID, A.locationFlag, A.locationID, A.typeID, A.itemID,
			T.typeName, IF(A.quantity, A.quantity, A.isSingleton) AS quantity, 
			count(*) - 1 AS subCount, buy, sell
			
		FROM evedata.assets A
		LEFT JOIN assets L ON A.itemID = L.locationID
		LEFT JOIN jitaPrice P ON A.typeID = P.itemID
		JOIN invTypes T ON A.typeID = T.typeID
		WHERE A.locationType != "other"
			AND A.characterID IN (SELECT tokenCharacterID FROM crestTokens WHERE characterID = ?)
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
			go getSubAssets(characterID, assets[index].ItemID, &assets[index].SubItems, errc, limit)
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

func getSubAssets(characterID int64, itemID int64, assets *[]Assets, errc chan error, limit chan bool) {
	limit <- true
	defer func() { <-limit }()
	characterID = characterID
	if err := database.Select(assets, `
		SELECT A.characterID, A.locationFlag, A.locationID, A.typeID, A.itemID,
			T.typeName, IF(A.quantity, A.quantity, A.isSingleton) AS quantity, 
			count(*) - 1 AS subCount, buy, sell
			
		FROM evedata.assets A
		LEFT JOIN evedata.assets L ON A.itemID = L.locationID
		LEFT JOIN jitaPrice P ON A.typeID = P.itemID
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
			go getSubAssets(characterID, a[index].ItemID, &a[index].SubItems, errc, limit)
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
