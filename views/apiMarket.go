package views

import (
	"encoding/json"
	"errors"
	"evedata/server"
	"net/http"
)

func init() {
	evedata.AddRoute(evedata.Route{"searchItems", "GET", "/J/searchItems", searchitemsPage})
}

type marketItemList struct {
	TypeID     int64  `db:"typeID"`
	TypeName   string `db:"typeName"`
	Categories string `db:"Categories"`
	Count      int64
}

// ARows bridge for old version
type ARows struct {
	Rows *[]marketItemList `json:"rows"`
}

func searchitemsPage(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	var q string
	q = r.FormValue("q")

	if len(q) < 2 {
		return http.StatusInternalServerError, errors.New("Query too short")
	}

	mIL := []marketItemList{}

	err := c.Db.Select(&mIL, `SELECT  T.typeID, typeName, CONCAT_WS(',', G5.marketGroupName, G4.marketGroupName, G3.marketGroupName, G2.marketGroupName, G.marketGroupName) AS Categories, count(*) AS count
           FROM market M
           INNER JOIN invTypes T ON M.typeID = T.typeID
           LEFT JOIN invMarketGroups G on T.marketGroupID = G.marketGroupID
           LEFT JOIN invMarketGroups G2 on G.parentGroupID = G2.marketGroupID
           LEFT JOIN invMarketGroups G3 on G2.parentGroupID = G3.marketGroupID
           LEFT JOIN invMarketGroups G4 on G3.parentGroupID = G4.marketGroupID
           LEFT JOIN invMarketGroups G5 on G4.parentGroupID = G5.marketGroupID

           WHERE done=0 AND T.marketGroupID IS NOT NULL AND typeName LIKE ?
           GROUP BY T.typeID
           ORDER BY typeName
           LIMIT 25`, "%"+q+"%")

	if err != nil {
		return http.StatusInternalServerError, err
	}

	var mRows ARows

	mRows.Rows = &mIL

	encoder := json.NewEncoder(w)
	encoder.Encode(mRows)

	return 200, nil
}
