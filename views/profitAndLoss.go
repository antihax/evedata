package views

import (
	"encoding/json"
	"errors"
	"strconv"

	"html/template"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
	"github.com/antihax/evedata/server"
	"github.com/antihax/evedata/templates"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("profitandloss", "GET", "/profitAndLoss", profitAndLossPage)
	evedata.AddAuthRoute("profitandloss", "GET", "/U/walletSummary", walletSummaryAPI)
}

func profitAndLossPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 60*60)
	p := newPage(r, "Profit and Loss Statement")
	templates.Templates = template.Must(template.ParseFiles("templates/profitAndLoss.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func walletSummaryAPI(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	var (
		err    error
		rangeI int64
	)

	setCache(w, 5*60)

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int64)
	if !ok {
		return http.StatusUnauthorized, errors.New("Unauthorized: Please log in.")
	}

	// Get range in days
	rangeTxt := r.FormValue("range")
	if rangeTxt != "" {
		rangeI, err = strconv.ParseInt(rangeTxt, 10, 64)
		if err != nil {
			return http.StatusNotFound, errors.New("Invalid range")
		}
	} else {
		return http.StatusInternalServerError, errors.New("range is required")
	}

	summary, err := models.GetWalletSummary(characterID, rangeI)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(summary)

	return 200, nil
}
