package evedata

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

/******************************************************************************
 * marketRegions JSON query
 *****************************************************************************/
type marketRegion struct {
	RegionID   int64  `db:"regionID"    json:"regionID"`
	Count      int64  `json:"count"`
	RegionName string `db:"regionName"  json:"regionName"`
}

// MarketRegions Query market regions from the database and return JSON to the
// user

func MarketRegions(c *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	mR := []marketRegion{}

	err := c.Db.Select(&mR, "SELECT regionID, regionName, count FROM tradeRegions GROUP BY regionID ORDER BY regionName")
	if err != nil {
		return 500, err
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	encoder := json.NewEncoder(w)
	encoder.Encode(mR)
	return 200, nil
}

/******************************************************************************
 * marketSellRegionItems JSON query
 *****************************************************************************/
type marketItems struct {
	StationName string `db:"stationName" json:"stationName"`
	StationID   string `db:"stationID"   json:"stationID"   `
	Quantity    string `db:"quantity"    json:"quantity"   `
	Price       string `db:"price"       json:"price"      `
}

// Rows is a list of rows for JSON conversion
type Rows struct {
	Rows *[]marketItems `json:"rows"`
}

// MarketSellRegionItems Query market sell orders for a user specified
// regionID and itemID query string and return JSON to the user
func MarketSellRegionItems(c *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var (
		regionID int
		itemID   int
		mRows    Rows
		err      error
	)

	mR := []marketItems{}

	regionID, err = strconv.Atoi(r.FormValue("regionID"))
	if err != nil {
		return 500, err
	}

	itemID, err = strconv.Atoi(r.FormValue("itemID"))
	if err != nil {
		return 500, err
	}

	err = c.Db.Select(&mR, `SELECT  remainingVolume AS quantity, price, stationName, M.stationID
                             FROM    market M
                             INNER JOIN staStations S ON S.stationID=M.stationID
                             WHERE      done=0 AND
                                      bid=0 AND
                                      M.regionID = ?
                                      AND typeID = ?
                             ORDER BY price ASC
                             `, regionID, itemID)

	if err != nil {
		return 500, err
	}

	mRows.Rows = &mR

	encoder := json.NewEncoder(w)
	encoder.Encode(mRows)
	return 200, nil
}

// MarketBuyRegionItems Query market buy orders for a user specified
// regionID and itemID query string and return JSON to the user
func MarketBuyRegionItems(c *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	var (
		regionID int
		itemID   int
		mRows    Rows
		err      error
	)

	mR := []marketItems{}

	regionID, err = strconv.Atoi(r.FormValue("regionID"))
	if err != nil {
		return 500, errors.New("Invalid regionID")
	}

	itemID, err = strconv.Atoi(r.FormValue("itemID"))
	if err != nil {
		return 500, errors.New("Invalid itemID")
	}

	err = c.Db.Select(&mR, `   SELECT  remainingVolume AS quantity, price, stationName, M.stationID
                             FROM    market M
                             INNER JOIN staStations S ON S.stationID=M.stationID
                             WHERE      done=0 AND
                                      bid=1 AND
                                      M.regionID = ?
                                      AND typeID = ?
                             ORDER BY price DESC
                             `, regionID, itemID)

	if err != nil {
		return 500, err
	}

	mRows.Rows = &mR

	encoder := json.NewEncoder(w)
	encoder.Encode(mRows)
	return 200, nil
}

/******************************************************************************
 * marketItemList JSON query
 *****************************************************************************/

type marketItemList struct {
	TypeID     int64  `db:"typeID"`
	TypeName   string `db:"typeName"`
	Categories string `db:"Categories"`
	Count      int64
}

type marketTree struct {
	Data      string        `json:"data"`
	Children  []*marketTree `json:"children,omitempty"`
	Attribute struct {
		ID int64 `json:"id,omitempty"`
	} `json:"attr,omitempty"`
}

// MarketItemLists queries the database for a user specified regionID
// returning a JSON list to the user.
func MarketItemLists(c *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	var id int
	var err error

	id, err = strconv.Atoi(r.FormValue("regionID"))
	if err != nil {
		return 500, err
	}

	Rows, err := c.Db.Queryx(`SELECT  T.typeID, typeName, CONCAT_WS(',', G5.marketGroupName, G4.marketGroupName, G3.marketGroupName, G2.marketGroupName, G.marketGroupName) AS Categories, count(*) AS count
           FROM    market M
           INNER JOIN invTypes T ON M.typeID = T.typeID
           LEFT JOIN invMarketGroups G on T.marketGroupID = G.marketGroupID
           LEFT JOIN invMarketGroups G2 on G.parentGroupID = G2.marketGroupID
           LEFT JOIN invMarketGroups G3 on G2.parentGroupID = G3.marketGroupID
           LEFT JOIN invMarketGroups G4 on G3.parentGroupID = G4.marketGroupID
           LEFT JOIN invMarketGroups G5 on G4.parentGroupID = G5.marketGroupID

           WHERE regionID = ? AND done=0 AND T.marketGroupID IS NOT NULL
           GROUP BY T.typeID
           ORDER BY Categories, typeName`, id)

	var (
		mTree       marketTree
		groups      []string
		last        [5]string
		lastNode    [6]*marketTree
		currentNode *marketTree
	)
	marketResult := marketItemList{}

	// Setup a root node for the Tree
	mTree.Data = "Root"
	lastNode[0] = &mTree

	// Since we used MySQL to order the results per group.
	// we can cheat here and build arrays at once time
	// without all that tedious searching.

	for Rows.Next() {
		err := Rows.StructScan(&marketResult)
		if err != nil {
			return 500, err
		}

		groups = strings.Split(marketResult.Categories, ",")

		// See if we changed MarketGroups
		for index, element := range groups {

			if element != last[index] {
				// We changed, lets make a new node.
				last[index] = element
				newNode := marketTree{}
				newNode.Data = element

				lastNode[index].Children = append(lastNode[index].Children, &newNode)
				lastNode[index+1] = &newNode

				currentNode = &newNode
			}
		}

		marketNode := marketTree{}
		marketNode.Data = marketResult.TypeName + " (" + strconv.FormatInt(marketResult.Count, 10) + ")"
		marketNode.Attribute.ID = marketResult.TypeID

		currentNode.Children = append(currentNode.Children, &marketNode)
	}

	// Skip the root node and JSONify.
	encoder := json.NewEncoder(w)
	encoder.Encode(mTree.Children)
	return 200, nil
}
