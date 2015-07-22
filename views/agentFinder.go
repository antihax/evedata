package views

import (
	"encoding/json"
	"evedata/server"
	"evedata/templates"
	"html/template"
	"net/http"
)

func init() {
	evedata.AddRoute(evedata.Route{"agentFinder", "GET", "/agentFinder", agentFinder})
	evedata.AddRoute(evedata.Route{"agentFinder", "GET", "/J/knownSpaceSystems", knownSpaceSystems})
}

func agentFinder(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	p := Page{
		Title: "EVE Online Market Browser",
	}

	templates.Templates = template.Must(template.ParseFiles("templates/agentFinder.html", templates.LayoutPath))
	err := templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

type solarSystemList struct {
	SolarSystemID   int64  `db:"solarSystemID"`
	SolarSystemName string `db:"solarSystemName"`
	Security        string `db:"security"`
}

// ARows bridge for old version
type aRows struct {
	Rows *[]solarSystemList `json:"rows"`
}

func knownSpaceSystems(c *evedata.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	var q, h string
	var sec float64
	q = r.FormValue("q")
	h = r.FormValue("hisec")
	if h == "true" {
		sec = 0.499999
	} else {
		sec = -1
	}

	sSL := []solarSystemList{}

	err := c.Db.Select(&sSL, `
		SELECT 
			solarSystemID,
			solarSystemName,
    		round(security, 1) AS security
		FROM 
			eve.mapSolarSystems 
		WHERE 
			regionID < 10999999 AND
			solarSystemName LIKE ? AND security > ?
		ORDER BY solarSystemName
        LIMIT 50
        `, q+"%", sec)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	var mRows aRows

	mRows.Rows = &sSL

	encoder := json.NewEncoder(w)
	encoder.Encode(mRows)

	return 200, nil
}
