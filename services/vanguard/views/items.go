package views

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/evedata/internal/strip"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("items", "GET", "/item", itemPage)
	vanguard.AddRoute("items", "GET", "/J/marketHistory", marketHistory)
}

func itemPage(w http.ResponseWriter, r *http.Request) {
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

	renderTemplate(w, "items.html", time.Hour*24*31, p)
}

func marketHistory(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		log.Println(err)
		return
	}

	v, err := models.GetMarketHistory(itemID, (int32)(regionID))
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour*12)
}
