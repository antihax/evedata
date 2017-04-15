package views

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/antihax/evedata/evedata"
	"github.com/antihax/evedata/models"
)

func init() {
	// Add routes to the http router
	evedata.AddRoute("searchItems", "GET", "/J/search", searchAPI)
}

// searchAPI for characters, alliances, corporations, and items.
func searchAPI(w http.ResponseWriter, r *http.Request) {

	// Get the query
	q := r.FormValue("q")
	q = strings.TrimSpace(q)

	// Make sure the query is at least three characters
	if len(q) <= 3 {
		httpErr(w, errors.New("Query too short"))
		return
	}

	// Do the search
	list, err := models.SearchNames(q)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Return the JSON representation
	json.NewEncoder(w).Encode(list)
}
