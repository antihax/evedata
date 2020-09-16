package models

import "github.com/guregu/null"

type ItemType struct {
	TypeID      int64   `db:"typeID" json:"typeID"`
	TypeName    string  `db:"typeName" json:"typeName"`
	Description string  `db:"description" json:"description"`
	Mass        float64 `db:"mass" json:"mass"`
	Volume      float64 `db:"volume" json:"volume"`
	Capacity    int64   `db:"capacity" json:"capacity"`
	PortionSize int64   `db:"portionSize" json:"portionSize"`
	RaceID      int64   `db:"raceID" json:"raceID"`
	RaceName    string  `db:"raceName" json:"raceName"`
	BasePrice   float64 `db:"basePrice" json:"basePrice"`
	IconID      int64   `db:"iconID" json:"iconID"`
	SoundID     int64   `db:"soundID" json:"soundID"`
	GraphicID   int64   `db:"graphicID" json:"graphicID"`
}

// Obtain Item information by ID.

func GetItem(id int64) (*ItemType, error) {
	ref := ItemType{}
	if err := database.QueryRowx(`
		SELECT 
			typeID, 
		    typeName, 
		    T.description, 
		    mass, 
		    volume, 
		    capacity, 
		    portionSize, 
		    IFNULL(T.raceID, 0) AS raceID,
		    IFNULL(raceName, "None") AS raceName,
		    IFNULL(basePrice, 0) AS basePrice, 
		    IFNULL(T.iconID, 0) AS iconID, 
		    IFNULL(soundID, 0) AS soundID,
		    graphicID
		 FROM invTypes T
		 LEFT OUTER JOIN chrRaces R ON T.raceID = R.raceID
		 WHERE typeID = ? LIMIT 1
			`, id).StructScan(&ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

type ItemAttributes struct {
	AttributeID         int64       `db:"attributeID" json:"attributeID"`
	Value               float64     `db:"value" json:"value"`
	AttributeName       string      `db:"attributeName" json:"attributeName"`
	Description         null.String `db:"description" json:"description"`
	CategoryID          int64       `db:"categoryID" json:"categoryID"`
	CategoryName        null.String `db:"categoryName" json:"categoryName"`
	CategoryDescription null.String `db:"categoryDescription" json:"categoryDescription"`
}

// Obtain Item Attributes by ID.

func GetItemAttributes(id int64) (*[]ItemAttributes, error) {
	ref := &[]ItemAttributes{}
	if err := database.Select(ref, `
		SELECT 
			A.attributeID, 
			IFNULL(valueFloat, valueInt) AS value,
		    IFNULL(displayName, IFNULL(attributeName, "UNKNOWN")) AS attributeName,
		    description,
		    T.categoryID,
		    categoryName,
		    C.categoryDescription
		
		FROM eve.dgmTypeAttributes A
		INNER JOIN dgmAttributeTypes T ON T.attributeID = A.attributeID
		INNER JOIN dgmAttributeCategories C ON C.categoryID = T.categoryID
		WHERE typeID = ?
		ORDER BY categoryID, attributeID
			`, id); err != nil {
		return nil, err
	}
	return ref, nil
}
