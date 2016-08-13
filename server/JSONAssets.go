package evedata

import (
	"encoding/json"

	"github.com/gorilla/sessions"

	"evedata/appContext"
	"evedata/models"
	"evedata/null"
	"net/http"
	"strconv"
	"time"
)

func init() {
	AddRoute(Route{"nextAssetCheck", "GET", "/U/nextAssetCheck", NextAssetCheck})
	AddRoute(Route{"assetStations", "GET", "/U/assetStations", AssetStations})
	AddRoute(Route{"assets", "GET", "/U/assets", Assets})
	AddRoute(Route{"assetCharacters", "GET", "/U/assetCharacters", AssetCharacters})
}

// NextAssetCheck ()
func NextAssetCheck(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {

	var (
		err         error
		characterID int

		nextUpdate struct {
			NextUpdate string `json:"nextUpdate"`
		}
	)

	characterID, err = strconv.Atoi(r.FormValue("characterID"))
	if err != nil {
		return 500, err
	}

	user := models.GetUser(r)

	if characterID != 0 && user != nil {
		var nextTime time.Time
		err = c.Db.Get(&nextTime, `
		   SELECT nextAssetCheck AS nextUpdate FROM characters
		   WHERE uid=? AND characterID = ?`, user.UID, characterID)

		nextUpdate.NextUpdate = nextTime.Format("Jan 2, 2006; 3:04pm (MST)")
	}

	if err != nil {
		return 500, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(nextUpdate)

	return 200, nil
}

// AssetCharacters ()
func AssetCharacters(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {

	type characterRow struct {
		Name        string `db:"name" json:"name"`
		CharacterID int64  `db:"characterID" json:"characterID"`
	}

	var (
		err error
	)

	rows := []characterRow{}

	user := models.GetUser(r)

	if user != nil {
		err = c.Db.Select(&rows, `
		   SELECT CONCAT(name, ' (', CAST(format(SUM(sell * A.quantity), 2) AS CHAR), ')') as name, C.characterID FROM characters C
		      INNER JOIN assets A ON A.characterID = C.characterID
		      INNER JOIN jitaPrice P ON A.typeID = P.itemID
		      WHERE C.uid = ?
		      GROUP BY name
		      ORDER BY SUM(sell * A.quantity) DESC
		      `, user.UID)
	}

	if err != nil {
		return 500, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(rows)

	return 200, nil
}

// AssetStations ()
func AssetStations(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {

	type stationRow struct {
		Name       string      `db:"name" json:"name"`
		Isk        null.String `db:"ISK" json:"ISK"`
		LocationID int64       `db:"locationID" json:"locationID"`
	}

	var (
		err         error
		characterID int
	)

	characterID, err = strconv.Atoi(r.FormValue("characterID"))
	if err != nil {
		return 500, err
	}

	rows := []stationRow{}

	user := models.GetUser(r)

	if characterID != 0 && user != nil {
		err = c.Db.Select(&rows, `
		   SELECT locationID, St.stationName AS name, sum(assets.quantity * P.buy) AS ISK
		   FROM assets
		   LEFT OUTER JOIN staStations AS St ON assets.locationID = St.stationID
		   LEFT OUTER JOIN mapSolarSystems AS S ON assets.locationID = S.solarSystemID
		   LEFT OUTER JOIN jitaPrice AS P ON assets.typeID = P.itemID
		   WHERE uid=? AND characterID=?
		   GROUP BY locationID HAVING name IS NOT NULL ORDER BY ISK DESC
		   `, user.UID, characterID)
	}

	if err != nil {
		return 500, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(rows)

	return 200, nil
}

// Assets ()
func Assets(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {

	type assetRow struct {
		Item      string `db:"typeName" json:"item"`
		ID        string `db:"typeID" json:"id"`
		Sell      string `db:"sell" json:"sell"`
		Buy       string `db:"buy" json:"buy"`
		BuyPrice  string `db:"buyPrice" json:"buyPrice"`
		SellPrice string `db:"sellPrice" json:"sellPrice"`
		Quantity  int64  `db:"quantity" json:"quantity"`
	}

	type aRows struct {
		Rows *[]assetRow `json:"rows"`
	}

	var (
		err         error
		characterID int
		locationID  int
		mRows       aRows
	)

	characterID, err = strconv.Atoi(r.FormValue("characterID"))
	if err != nil {
		return 500, err
	}

	locationID, err = strconv.Atoi(r.FormValue("locationID"))
	if err != nil {
		return 500, err
	}

	rows := []assetRow{}

	user := models.GetUser(r)

	if characterID != 0 && user != nil {
		err = c.Db.Select(&rows, `
		   SELECT typeName, T.typeID, A.quantity, buy, sell, buy * A.quantity AS buyPrice, sell * A.quantity AS sellPrice  FROM assets A
		   INNER JOIN invTypes T ON A.typeID = T.typeID
		   INNER JOIN jitaPrice P ON A.typeID = P.itemID
		   WHERE uid = ? AND characterID = ? AND locationID = ?
		   `, user.UID, characterID, locationID)
	}

	if err != nil {
		return 500, err
	}

	mRows.Rows = &rows
	encoder := json.NewEncoder(w)
	encoder.Encode(mRows)

	return 200, nil
}
