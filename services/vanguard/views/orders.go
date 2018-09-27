package views

import (
	"errors"
	"log"
	"strconv"
	"time"

	"net/http"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"

	"github.com/antihax/goesi"
)

func init() {
	vanguard.AddRoute("GET", "/orders", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "orders.html", time.Hour*24*31, newPage(r, "Order Information"))
	})

	vanguard.AddAuthRoute("GET", "/U/orders", ordersAPI)
	vanguard.AddAuthRoute("GET", "/U/orderCharacters", orderCharactersAPI)
}

func orderCharactersAPI(w http.ResponseWriter, r *http.Request) {
	var err error
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	ch, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		log.Println(err)
		httpErrCode(w, errors.New("could not find character ID for orders"), http.StatusUnauthorized)
		return
	}

	v, err := models.GetOrderCharacters(ch.CharacterID, ch.CharacterOwnerHash)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	if len(v) == 0 {
		httpErrCode(w, errors.New("No order characters"), http.StatusNotFound)
		return
	}

	renderJSON(w, v, time.Minute)
}

func ordersAPI(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		filterCharacterID int32
	)

	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	ch, ok := s.Values["character"].(goesi.VerifyResponse)
	if !ok {
		httpErrCode(w, errors.New("could not find character ID for order API"), http.StatusUnauthorized)
		return
	}

	// Get arguments
	filter := r.FormValue("filterCharacterID")
	if filter != "" {
		filterCharacterID64, err := strconv.ParseInt(filter, 10, 64)
		if err != nil {
			log.Println(err)
			httpErrCode(w, err, http.StatusNotFound)
			return
		}
		filterCharacterID = int32(filterCharacterID64)
	}

	v, err := models.GetOrders(ch.CharacterID, ch.CharacterOwnerHash, filterCharacterID)
	if err != nil {
		log.Println(err)
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Minute)
}
