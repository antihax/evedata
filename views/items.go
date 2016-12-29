package views

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/strip"
	"github.com/antihax/evedata/templates"

	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("items", "GET", "/item", itemPage)
	evedata.AddRoute("items", "GET", "/J/marketHistory", marketHistory)
}

func itemPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "Unknown Item")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid item ID. Please provide an ?id=")
	}

	errc := make(chan error)

	// Get the item information
	go func() {
		ref, err := models.GetItem(id)
		if err != nil {
			errc <- err
			return
		}
		ref.Description = strip.StripTags(ref.Description)
		p["Item"] = ref
		p["Title"] = ref.TypeName
		errc <- nil
	}()
	// Get the item information
	go func() {
		ref, err := models.GetItemAttributes(id)
		if err != nil {
			errc <- err
			return
		}
		p["ItemAttributes"] = ref
		errc <- nil
	}()
	// clear the error channel
	for i := 0; i < 2; i++ {
		err := <-errc
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	templates.Templates = template.Must(template.ParseFiles("templates/items.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func marketHistory(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*240)
	region := r.FormValue("regionID")
	item := r.FormValue("itemID")

	itemID, err := strconv.ParseInt(item, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}

	regionID, err := strconv.Atoi(region)

	v, err := models.GetMarketHistory(itemID, (int32)(regionID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)

	return 200, nil
}
