package models

import "evedata/null"

type IskPerLP struct {
	ItemName     string      `db:"itemName" json:"itemName"`
	TypeName     string      `db:"typeName" json:"typeName"`
	JitaPrice    float64     `db:"JitaPrice" json:"jitaPrice"`
	ItemCost     float64     `db:"itemCost" json:"itemCost"`
	IskPerLP     int64       `db:"iskPerLP" json:"iskPerLP"`
	JitaVolume   int64       `db:"JitaVolume" json:"jitaVolume"`
	Requirements null.String `db:"requirements" json:"requirements"`
}

func GetISKPerLP(corporationName string) ([]IskPerLP, error) {
	s := []IskPerLP{}
	if err := database.Select(&s, `
		SELECT itemName, Lp.typeName, JitaPrice, itemCost, iskPerLP, JitaVolume, GROUP_CONCAT(quantity, " x ", T.typeName SEPARATOR '<br>\n') AS requirements
			FROM iskPerLp Lp
			LEFT JOIN lpOfferRequirements R ON Lp.offerID = R.offerID
			LEFT JOIN invTypes T ON R.typeID = T.typeID
			WHERE itemName = ?
			GROUP BY Lp.typeName
			ORDER BY ISKperLP DESC;
	;`, corporationName); err != nil {

		return nil, err
	}
	return s, nil
}

type IskPerLPCorporation struct {
	ItemName string `db:"itemName" json:"itemName" `
}

func GetISKPerLPCorporations() ([]IskPerLPCorporation, error) {
	s := []IskPerLPCorporation{}
	if err := database.Select(&s, `SELECT DISTINCT itemName FROM iskPerLp ORDER BY itemName ASC;`); err != nil {
		return nil, err
	}
	return s, nil
}

type ArbitrageCalculatorStations struct {
	StationName string `db:"stationName" json:"stationName" `
	StationID   string `db:"stationID" json:"stationID" `
}

func GetArbitrageCalculatorStations() ([]ArbitrageCalculatorStations, error) {
	s := []ArbitrageCalculatorStations{}
	if err := database.Select(&s, `
		SELECT stationID, stationName
			FROM    marketStations
			WHERE 	Count > 4000
			ORDER BY stationName
	`); err != nil {
		return nil, err
	}
	return s, nil
}

type ArbitrageCalculator struct {
	TypeID      int64  `db:"typeID" json:"typeID" `
	TypeName    string `db:"typeName" json:"typeName" `
	PortionSize int64  `db:"portionSize" json:"portionSize" `
	Buys        string `db:"buys" json:"buys" `
	Vol         int64  `db:"vol" json:"vol" `
	Price       int64  `db:"price" json:"price" `
}

func GetArbitrageCalculator(hours int64, stationID int64, minVolume int64, maxPrice int64) ([]ArbitrageCalculator, error) {
	s := []ArbitrageCalculator{}
	if err := database.Select(&s, `
		SELECT  market.typeID AS typeID, typeName, portionSize, count(*) as buys, market_vol.quantity / 2 as vol, max(price) AS price
		FROM    market, invTypes, market_vol
		WHERE   market.done = 0 AND
		        market.typeID = market_vol.itemID AND
		        market.regionID = market_vol.regionID AND
		        market.typeID = invTypes.typeID AND
		        reported >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL ? DAY_HOUR) AND
		        market.stationID = ?
		        AND bid = 1
		        AND done = 0 AND
		        market_vol.quantity /2 > ?
		GROUP BY market.typeID
		HAVING price < ?`); err != nil {
		return nil, err
	}
	return s, nil
}
