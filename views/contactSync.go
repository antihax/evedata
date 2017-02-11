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
	"github.com/antihax/evedata/templates"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("ContactSync", "GET", "/contactSync", contactSyncPage)
	evedata.AddAuthRoute("ContactSync", "PUT", "/U/contactSync", apiAddContactSync)
	evedata.AddAuthRoute("ContactSync", "GET", "/U/contactSync", apiGetContactSyncs)
	evedata.AddAuthRoute("ContactSync", "DELETE", "/U/contactSync", apiDeleteContactSync)
}

func contactSyncPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	p := newPage(r, "Contact Copiers")
	templates.Templates = template.Must(template.ParseFiles("templates/contactSync.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func apiAddContactSync(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	type localContactSync struct {
		Source      int `json:",string"`
		Destination int `json:",string"`
	}
	var cc localContactSync

	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized: Please log in.")
	}

	if r.Body == nil {
		return http.StatusBadRequest, errors.New("No Data Received")
	}

	err := json.NewDecoder(r.Body).Decode(&cc)
	if err != nil {
		return http.StatusNotFound, err
	}

	if err := models.AddContactSync(characterID, cc.Source, cc.Destination); err != nil {
		return http.StatusConflict, err
	}

	return 200, nil
}

func apiGetContactSyncs(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized: Please log in.")
	}

	cc, err := models.GetContactSyncs(characterID)
	if err != nil {
		return http.StatusNotFound, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(cc)

	return 200, nil
}

func apiDeleteContactSync(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 0)

	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized: Please log in.")
	}

	destination, err := strconv.Atoi(r.FormValue("destination"))
	if err != nil {
		return http.StatusNotFound, errors.New("Invalid destination")
	}

	if err := models.DeleteContactSync(characterID, destination); err != nil {
		return http.StatusConflict, err
	}

	return 200, nil
}
