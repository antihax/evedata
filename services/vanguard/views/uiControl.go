package views

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/goesi"
)

func init() {
	vanguard.AddAuthRoute("POST", "/X/setDestination", setDestination)
	vanguard.AddAuthRoute("POST", "/X/addDestination", addDestination)
	vanguard.AddAuthRoute("POST", "/X/openMarketWindow", openMarketWindow)
}

func openMarketWindow(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())
	c := vanguard.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	// Get the typeID
	typeID, err := strconv.ParseInt(r.FormValue("typeID"), 10, 32)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Get the character
	tokenCharacterID, _ := strconv.ParseInt(r.FormValue("characterID"), 10, 32)
	if tokenCharacterID > 0 {
		tokenSource, err := c.TokenStore.GetTokenSource(characterID, int32(tokenCharacterID))
		if err != nil {
			log.Println(err)
			return
		}
		auth := context.WithValue(r.Context(), goesi.ContextOAuth2, tokenSource)
		openMarket(c, auth, int32(typeID))
	} else {
		// Get the control character authentication
		auth, err := getCursorCharacterAuth(c, s)
		if err != nil {
			httpErr(w, err)
			return
		}
		openMarket(c, auth, int32(typeID))
	}
}

func openMarket(c *vanguard.Vanguard, auth context.Context, typeID int32) error {
	_, err := c.ESI.ESI.UserInterfaceApi.PostUiOpenwindowMarketdetails(auth, typeID, nil)
	return err
}

func setDestination(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())
	c := vanguard.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	_, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	// Get the destinationID for the location
	destinationIDTxt := r.FormValue("destinationID")
	destinationID, err := strconv.ParseInt(destinationIDTxt, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Get the control character authentication
	auth, err := getCursorCharacterAuth(c, s)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Set the destination
	res, err := c.ESI.ESI.UserInterfaceApi.PostUiAutopilotWaypoint(auth, false, true, destinationID, nil)
	if err != nil {
		if res != nil {
			httpErrCode(w, err, res.StatusCode)
			return
		}
		httpErr(w, err)
		return
	}
}

func addDestination(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())
	c := vanguard.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	_, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	// Get the destinationID for the location
	destinationIDTxt := r.FormValue("destinationID")
	destinationID, err := strconv.ParseInt(destinationIDTxt, 10, 64)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Get the control character authentication
	auth, err := getCursorCharacterAuth(c, s)
	if err != nil {
		httpErr(w, err)
		return
	}

	// Set the destination
	res, err := c.ESI.ESI.UserInterfaceApi.PostUiAutopilotWaypoint(auth, false, false, destinationID, nil)
	if err != nil {
		if res != nil {
			httpErrCode(w, err, res.StatusCode)
			return
		}
		httpErr(w, err)
		return
	}
}
