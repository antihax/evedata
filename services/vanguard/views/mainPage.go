package views

import (
	"net/http"
	"time"

	"github.com/antihax/evedata/services/vanguard"
)

func init() {
	vanguard.AddRoute("GET", "/",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"mainPage.html",
				time.Hour*24*31,
				newPage(r, "EVE Online Intel Data"))
		})
}
