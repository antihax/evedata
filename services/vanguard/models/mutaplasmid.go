package models

import (
	"errors"
	"strings"
)

type MutaplasmidData struct {
	MetaData string `db:"metaData"`
	Data     string `db:"data"`
}

var MutaplasmidTypes = map[string]string{
	"100mn Afterburner":         `T.typeName LIKE "%100mn%" AND groupID = 46 `,
	"10mn Afterburner":          `T.typeName LIKE "%10mn%" AND groupID = 46 `,
	"1mn Afterburner":           `T.typeName LIKE "%1mn%" AND groupID = 46 `,
	"500mn Microwarpdrive":      `T.typeName LIKE "%500mn%" AND groupID = 46`,
	"50mn Microwarpdrive":       `T.typeName LIKE "%50mn%" AND groupID = 46`,
	"5mn Microwarpdrive":        `T.typeName LIKE "%5mn%" AND groupID = 46`,
	"1600mm Armor Plate":        `T.typeName LIKE "%1600mm%" AND groupID = 329`,
	"Heavy Energy Neutralizer":  `T.typeName LIKE "%heavy%" AND groupID = 71`,
	"Large Armor Repairer":      `T.typeName LIKE "%large%" AND groupID = 62`,
	"Large Shield Booster":      `T.typeName LIKE "%large%" AND groupID = 40 AND T.typeName NOT LIKE "%x-large%"`,
	"Large Shield Extender":     `T.typeName LIKE "%large%" AND groupID = 38`,
	"400-800mm Armor Plate":     `(T.typeName LIKE "%400mm%" OR T.typeName LIKE "%800mm%") AND groupID = 329`,
	"Medium Energy Neutralizer": `T.typeName LIKE "%medium%" AND groupID = 71`,
	"Medium Armor Repairer":     `T.typeName LIKE "%medium%" AND groupID = 62`,
	"Medium Shield Booster":     `T.typeName LIKE "%medium%" AND groupID = 40`,
	"Medium Shield Extender":    `T.typeName LIKE "%medium%" AND groupID = 38`,
	"100-200mm Armor Plate":     `(T.typeName LIKE "%100mm%" OR T.typeName LIKE "%200mm%") AND groupID = 329`,
	"Small Energy Neutralizer":  `T.typeName LIKE "%small%" AND groupID = 71`,
	"Small Armor Repairer":      `T.typeName LIKE "%small%" AND groupID = 62`,
	"Small Shield Booster":      `T.typeName LIKE "%small%" AND groupID = 40`,
	"Small Shield Extender":     `T.typeName LIKE "%small%" AND groupID = 38`,
	"Stasis Webifier":           `T.typeName LIKE "%stasis%"  AND groupID = 65 AND T.typeName NOT LIKE "%civilian%"`,
	"Warp Disruptor":            `T.typeName LIKE "%disruptor%" AND groupID = 52 AND T.typeName NOT LIKE "%heavy%" AND T.typeName NOT LIKE "%civilian%"`,
	"Warp Scrambler":            `T.typeName LIKE "%scrambler%" AND groupID = 52 AND T.typeName NOT LIKE "%heavy%" `,
	"X-Large Shield Booster":    `T.typeName LIKE "%x-large%" AND groupID = 40 AND T.typeName NOT LIKE "% large%"`,
}

// Obtain Item Attributes by ID.
// [BENCHMARK] 0.125 sec / 0.000 sec
func GetMutaplasmidData(mutaplasmidType string) (*MutaplasmidData, error) {
	type mutaplasmidData struct {
		MetaData string `db:"metaData"`
		Data     string `db:"data"`
	}

	mpt, ok := MutaplasmidTypes[mutaplasmidType]
	if !ok {
		return nil, errors.New("Unknown mutaplasmid type")
	}

	data := []mutaplasmidData{}
	if err := database.Select(&data, `
		SELECT 
			GROUP_CONCAT(
				DISTINCT CONCAT('{',
			'"key": "', attributeName, '", ',
			'"name": "', displayName, '"',
			"}") ORDER BY attributeName
			) AS metaData,
			CONCAT(
			'{', 
				'"typeName": "', typeName, '"' ', ',
				'"price": ', mean, ', ', 
				GROUP_CONCAT(
					DISTINCT CONCAT( '"', attributeName, '": ', IFNULL( valueInt, valueFloat ) ) ORDER BY attributeName ASC
				), '}') AS data
			FROM eve.invTypes T
			INNER JOIN dgmTypeAttributes A ON A.typeID = T.typeID  AND attributeID IN (50,6,73,84,105,20,72,983,796,54,30,68,97,554,1159)
			INNER JOIN dgmAttributeTypes AT ON AT.attributeID = A.attributeID
			INNER JOIN evedata.market_history H ON H.itemID = T.typeID  AND H.date >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY) AND H.regionID IN (10000002, 10000043, 10000030, 10000032, 30002053)
			WHERE `+mpt+`
			GROUP BY T.typeID, H.date, H.regionID
			ORDER BY mean ASC
			`); err != nil {
		return nil, err
	}
	strs := make([]string, len(data))
	for i, v := range data {
		strs[i] = v.Data
	}

	if len(data) == 0 {
		return nil, errors.New("no data found")
	}

	return &MutaplasmidData{
		MetaData: data[0].MetaData,
		Data:     strings.Join(strs, ","),
	}, nil
}

//			INNER JOIN dgmTypeAttributes A ON A.typeID = T.typeID  AND attributeID IN (50,6,73,84,105,20,72,983,796,54,30,68,97,554,1159)
