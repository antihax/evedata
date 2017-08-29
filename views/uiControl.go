package views

import (
	"net/http"
	"strconv"

	"github.com/antihax/evedata/evedata"
)

func init() {
	evedata.AddAuthRoute("ui-control", "POST", "/X/setDestination", setDestination)
	evedata.AddAuthRoute("ui-control", "POST", "/X/addDestination", addDestination)
	evedata.AddAuthRoute("ui-control", "POST", "/X/openMarketWindow", openMarketWindow)
}

func openMarketWindow(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())
	c := evedata.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	_, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
		return
	}

	// Get the destinationID for the location
	typeIDTxt := r.FormValue("typeID")
	typeID, err := strconv.ParseInt(typeIDTxt, 10, 32)
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
	res, err := c.ESI.ESI.UserInterfaceApi.PostUiOpenwindowMarketdetails(auth, (int32)(typeID), nil)
	if err != nil {
		if res != nil {
			httpErrCode(w, res.StatusCode)
			return
		}
		httpErr(w, err)
		return
	}
}

func setDestination(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())
	c := evedata.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	_, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
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
			httpErrCode(w, res.StatusCode)
			return
		}
		httpErr(w, err)
		return
	}
}

func addDestination(w http.ResponseWriter, r *http.Request) {
	setCache(w, 0)
	s := evedata.SessionFromContext(r.Context())
	c := evedata.GlobalsFromContext(r.Context())

	// Get the sessions main characterID
	_, ok := s.Values["characterID"].(int64)
	if !ok {
		httpErrCode(w, http.StatusUnauthorized)
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
			httpErrCode(w, res.StatusCode)
			return
		}
		httpErr(w, err)
		return
	}
}
