package evedata

import (
	"encoding/json"
	"evedata/appContext"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

func init() {
	AddRoute("marketRegions", "GET", "/J/marketRegions", MarketRegions)
	AddRoute("marketItemLists", "GET", "/J/marketItemLists", MarketItemLists)
}

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
func MarketRegions(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
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
func MarketItemLists(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {

	var regionID int
	var err error
	var Rows *sqlx.Rows
	regionID, err = strconv.Atoi(r.FormValue("regionID"))
	if err != nil {
		regionID = 0
	}

	if regionID == 0 {
		Rows, err = c.Db.Queryx(`SELECT  T.typeID, typeName, CONCAT_WS(',', G5.marketGroupName, G4.marketGroupName, G3.marketGroupName, G2.marketGroupName, G.marketGroupName) AS Categories, count(*) AS count
           FROM    market M
           INNER JOIN invTypes T ON M.typeID = T.typeID
           LEFT JOIN invMarketGroups G on T.marketGroupID = G.marketGroupID
           LEFT JOIN invMarketGroups G2 on G.parentGroupID = G2.marketGroupID
           LEFT JOIN invMarketGroups G3 on G2.parentGroupID = G3.marketGroupID
           LEFT JOIN invMarketGroups G4 on G3.parentGroupID = G4.marketGroupID
           LEFT JOIN invMarketGroups G5 on G4.parentGroupID = G5.marketGroupID

           WHERE done=0 AND T.marketGroupID IS NOT NULL
           GROUP BY T.typeID
           ORDER BY Categories, typeName`)
	} else {
		Rows, err = c.Db.Queryx(`SELECT  T.typeID, typeName, CONCAT_WS(',', G5.marketGroupName, G4.marketGroupName, G3.marketGroupName, G2.marketGroupName, G.marketGroupName) AS Categories, count(*) AS count
           FROM    market M
           INNER JOIN invTypes T ON M.typeID = T.typeID
           LEFT JOIN invMarketGroups G on T.marketGroupID = G.marketGroupID
           LEFT JOIN invMarketGroups G2 on G.parentGroupID = G2.marketGroupID
           LEFT JOIN invMarketGroups G3 on G2.parentGroupID = G3.marketGroupID
           LEFT JOIN invMarketGroups G4 on G3.parentGroupID = G4.marketGroupID
           LEFT JOIN invMarketGroups G5 on G4.parentGroupID = G5.marketGroupID

           WHERE regionID = ? AND done=0 AND T.marketGroupID IS NOT NULL
           GROUP BY T.typeID
           ORDER BY Categories, typeName`, regionID)
	}
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
