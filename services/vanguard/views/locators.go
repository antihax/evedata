package views

import (
	"net/http"
	"time"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("GET", "/locatorResponses",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"locatorResponses.html",
				time.Hour*24*31,
				newPage(r, "Locator Responses"))
		})
	vanguard.AddAuthRoute("GET", "/U/locatorResponses", apiGetLocatorResponses)
}

func apiGetLocatorResponses(w http.ResponseWriter, r *http.Request) {
	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	v, err := models.GetLocatorResponses(characterID)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Minute)

	if err = s.Save(r, w); err != nil {
		httpErr(w, err)
		return
	}
}
