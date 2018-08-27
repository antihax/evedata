package views

import (
	"net/http"
	"time"

	"github.com/antihax/evedata/services/vanguard"
)

func init() {
	vanguard.AddRoute("GET", "/killmail", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "killmail.html", time.Hour*24*31, newPage(r, "Killmail"))
	})
}
