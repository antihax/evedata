package views

import (
	"errors"
	"evedata/appContext"
	"evedata/models"
	"evedata/server"
	"evedata/strip"
	"evedata/templates"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("entity", "GET", "/alliance", alliancePage)
	evedata.AddRoute("entity", "GET", "/corporation", corporationPage)
	evedata.AddRoute("entity", "GET", "/character", characterPage)
}

func alliancePage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "Unknown Alliance")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid ID. Please provide an ?id=")
	}

	errc := make(chan error)

	// Get known wars. This takes the longest.
	go func() {
		ref, err := models.GetWarsForEntityByID(id)
		p["Wars"] = ref
		errc <- err
	}()

	// Get the alliance information
	go func() {
		ref, err := models.GetAlliance(id)
		ref.Description = strip.StripTags(ref.Description)
		p["Alliance"] = ref
		p["Title"] = ref.AllianceName
		errc <- err
	}()

	// Get the alliance members
	go func() {
		ref, err := models.GetAllianceMembers(id)
		p["AllianceMembers"] = ref
		errc <- err
	}()

	// clear the error channel
	for i := 0; i < 3; i++ {
		err := <-errc
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	templates.Templates = template.Must(template.ParseFiles("templates/entities.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func corporationPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "Unknown Corporation")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid ID. Please provide an ?id=")
	}

	errc := make(chan error)

	// Get known wars. This takes the longest.
	go func() {
		ref, err := models.GetWarsForEntityByID(id)
		p["Wars"] = ref
		errc <- err
	}()

	// Get the alliance information
	go func() {
		ref, err := models.GetCorporation(id)
		p["Corporation"] = ref
		p["Title"] = ref.CorporationName
		errc <- err
	}()

	// clear the error channel
	for i := 0; i < 2; i++ {
		err := <-errc
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	templates.Templates = template.Must(template.ParseFiles("templates/entities.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func characterPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "Unknown Character")

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Invalid ID. Please provide an ?id=")
	}

	errc := make(chan error)

	// Get the character information
	go func() {
		ref, err := models.GetCharacter(id)

		p["Character"] = ref
		p["Title"] = ref.CharacterName
		errc <- err
	}()

	// clear the error channel
	for i := 0; i < 1; i++ {
		err := <-errc
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	templates.Templates = template.Must(template.ParseFiles("templates/entities.html", templates.LayoutPath))
	err = templates.Templates.ExecuteTemplate(w, "base", p)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
