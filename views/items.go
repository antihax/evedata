package views

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/strip"
	"github.com/antihax/evedata/templates"
)

func init() {
	evedata.AddRoute("items", "GET", "/item", itemPage)
	evedata.AddRoute("items", "GET", "/J/marketHistory", marketHistory)
}

func itemPage(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60)
	p := newPage(r, "Unknown Item")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
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
			httpErr(w, err)
			return
		}
	}

	templates.Templates = template.Must(template.ParseFiles("templates/items.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)
	if err != nil {
		httpErr(w, err)
		return
	}
}

func marketHistory(w http.ResponseWriter, r *http.Request) {
	setCache(w, 60*60*24)
	region := r.FormValue("regionID")
	item := r.FormValue("itemID")

	itemID, err := strconv.ParseInt(item, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	if err != nil {
		httpErr(w, err)
		return
	}

	regionID, err := strconv.Atoi(region)

	v, err := models.GetMarketHistory(itemID, (int32)(regionID))
	if err != nil {
		httpErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}
