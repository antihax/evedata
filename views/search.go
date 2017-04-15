package views

import (
	"errors"
	"net/http"

	"github.com/antihax/evedata/evedata"
)

func init() {
	evedata.AddRoute("searchItems", "GET", "/J/searchItems", searchAPI)
}

func searchAPI(w http.ResponseWriter, r *http.Request) {
	var q string
	q = r.FormValue("q")

	if len(q) < 2 {
		httpErr(w, errors.New("Query too short"))
		return
	}

	//json.NewEncoder(w).Encode(mRows)
}
