package models

import "github.com/guregu/null"

func AddLPOffer(offerID int64, corporationID int64, typeID int64, quantity int64, lpCost int64, akCost, iskCost int64) error {
	if _, err := database.Exec(`
		INSERT IGNORE INTO evedata.lpOffers
			(offerID,corporationID,typeID,quantity,lpCost,akCost,iskCost)
			VALUES(?,?,?,?,?,?,?);
	`, offerID, corporationID, typeID, quantity, lpCost, akCost, iskCost); err != nil {
		return err
	}
	return nil
}

func AddLPOfferRequirements(offerID int64, typeID int64, quantity int64) error {
	if _, err := database.Exec(`INSERT IGNORE INTO evedata.lpOfferRequirements (offerID,typeID,quantity) VALUES(?,?,?);`,
		offerID, typeID, quantity); err != nil {
		return err
	}
	return nil
}

type IskPerLP struct {
	ItemName      string      `db:"itemName" json:"itemName"`
	TypeID        int64       `db:"typeID" json:"typeID"`
	TypeName      string      `db:"typeName" json:"typeName"`
	JitaPrice     float64     `db:"JitaPrice" json:"jitaPrice"`
	ItemCost      float64     `db:"itemCost" json:"itemCost"`
	IskPerLP      int64       `db:"iskPerLP" json:"iskPerLP"`
	JitaVolume    int64       `db:"JitaVolume" json:"jitaVolume"`
	IskVolume     float64     `db:"iskVolume" json:"iskVolume"`
	Requirements  null.String `db:"requirements" json:"requirements"`
	AvailableFrom string      `db:"availableFrom" json:"availableFrom"`
}

func GetISKPerLPByConversion() ([]IskPerLP, error) {
	s := []IskPerLP{}
	if err := database.Select(&s, `
		SELECT DISTINCT typeID, typeName, JitaPrice, itemCost, iskPerLP, JitaVolume, JitaVolume*JitaPrice AS iskVolume, GROUP_CONCAT(itemName ORDER BY itemName ASC SEPARATOR ', ') AS availableFrom
			FROM evedata.iskPerLp Lp
			GROUP BY typeID, lpCost
			ORDER BY ISKperLP DESC
			LIMIT 500
	;`); err != nil {

		return nil, err
	}
	return s, nil
}

func GetISKPerLP(corporationName string) ([]IskPerLP, error) {
	s := []IskPerLP{}
	if err := database.Select(&s, `
		SELECT itemName, Lp.typeID, Lp.typeName, JitaPrice, itemCost, iskPerLP, JitaVolume, JitaVolume*JitaPrice AS iskVolume, GROUP_CONCAT(quantity, " x ", T.typeName SEPARATOR '<br>\n') AS requirements
			FROM evedata.iskPerLp Lp
			LEFT JOIN evedata.lpOfferRequirements R ON Lp.offerID = R.offerID
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
	if err := database.Select(&s, `SELECT DISTINCT itemName FROM evedata.iskPerLp ORDER BY itemName ASC;`); err != nil {
		return nil, err
	}
	return s, nil
}
