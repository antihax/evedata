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
	evedata.AddRoute("searchItems", "GET", "/J/searchEntities", searchEntitiesAPI)
	evedata.AddRoute("searchItems", "GET", "/search", searchRouter)
}

// searchAPI for characters, alliances, corporations, and items.
func searchEntitiesAPI(w http.ResponseWriter, r *http.Request) {
	setCache(w, 12*60*60)
	// Get the query
	q := r.FormValue("q")
	q = strings.TrimSpace(q)

	// Make sure the query is at least three characters
	if len(q) <= 3 {
		httpErr(w, errors.New("Query too short"))
		return
	}

	// Do the search
	list, err := models.SearchEntities(q)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Return the JSON representation
	json.NewEncoder(w).Encode(list)
}

// searchAPI for characters, alliances, corporations, and items.
func searchAPI(w http.ResponseWriter, r *http.Request) {
	setCache(w, 12*60*60)
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

// searchAPI for characters, alliances, corporations, and items.
func searchRouter(w http.ResponseWriter, r *http.Request) {
	var endPoint string

	id := r.FormValue("id")
	entityType := strings.ToLower(r.FormValue("type"))

	switch entityType {
	case "character":
		endPoint = "/character?id=" + id
	case "alliance":
		endPoint = "/alliance?id=" + id
	case "corporation":
		endPoint = "/corporation?id=" + id
	case "item":
		endPoint = "/item?id=" + id
	default:
		httpErr(w, errors.New("Unknown endpoint"))
		return
	}

	http.Redirect(w, r, endPoint, 302)
	httpErrCode(w, nil, http.StatusMovedPermanently)
}
