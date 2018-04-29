package views

import (
	"strconv"
	"time"

	"net/http"

	"github.com/antihax/evedata/services/vanguard"
	"github.com/antihax/evedata/services/vanguard/models"
)

func init() {
	vanguard.AddRoute("profitandloss", "GET", "/profitAndLoss",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"profitAndLoss.html",
				time.Hour*24*31,
				newPage(r, "Profit and Loss Statement"))
		})
	vanguard.AddAuthRoute("profitandloss", "GET", "/U/walletSummary", walletSummaryAPI)
}

func walletSummaryAPI(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		rangeI int64
	)

	s := vanguard.SessionFromContext(r.Context())

	// Get the sessions main characterID
	characterID, ok := s.Values["characterID"].(int32)
	if !ok {
		httpErrCode(w, nil, http.StatusUnauthorized)
		return
	}

	// Get range in days
	rangeTxt := r.FormValue("range")
	if rangeTxt != "" {
		rangeI, err = strconv.ParseInt(rangeTxt, 10, 64)
		if err != nil {
			httpErrCode(w, err, http.StatusBadRequest)
			return
		}
	} else {
		httpErrCode(w, err, http.StatusBadRequest)
		return
	}

	v, err := models.GetWalletSummary(characterID, rangeI)
	if err != nil {
		httpErr(w, err)
		return
	}

	renderJSON(w, v, time.Hour)
}
