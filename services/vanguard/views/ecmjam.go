package views

import (
	"net/http"
	"time"

	"github.com/antihax/evedata/services/vanguard"
)

func init() {
	vanguard.AddRoute("ecmjam", "GET", "/ecmjam",
		func(w http.ResponseWriter, r *http.Request) {
			renderTemplate(w,
				"ecmjam.html",
				time.Hour*24*31,
				newPage(r, "EVE ECM Jam"))
		})
}
