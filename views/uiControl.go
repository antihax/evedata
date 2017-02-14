package views

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/server"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddAuthRoute("ui-control", "POST", "/X/setDestination", setDestination)
	evedata.AddAuthRoute("ui-control", "POST", "/X/addDestination", addDestination)
}

func setDestination(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)

	// Get the sessions main characterID
	_, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized: Please log in.")
	}

	// Get the destinationID for the location
	destinationIDTxt := r.FormValue("destinationID")
	destinationID, err := strconv.ParseInt(destinationIDTxt, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Get the control character authentication
	auth, err := getCursorCharacterAuth(c, s)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Set the destination
	res, err := c.ESI.V2.UserInterfaceApi.PostUiAutopilotWaypoint(auth, destinationID, true, false, nil)
	if err != nil {
		if res != nil {
			return res.StatusCode, err
		}
		return http.StatusInternalServerError, err
	}

	// Return the status code from CCP.
	return res.StatusCode, nil
}

func addDestination(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)

	// Get the sessions main characterID
	_, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized: Please log in.")
	}

	// Get the destinationID for the location
	destinationIDTxt := r.FormValue("destinationID")
	destinationID, err := strconv.ParseInt(destinationIDTxt, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Get the control character authentication
	auth, err := getCursorCharacterAuth(c, s)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Set the destination
	res, err := c.ESI.V2.UserInterfaceApi.PostUiAutopilotWaypoint(auth, destinationID, false, false, nil)
	if err != nil {
		if res != nil {
			return res.StatusCode, err
		}
		return http.StatusInternalServerError, err
	}

	// Return the status code from CCP.
	return res.StatusCode, nil
}
