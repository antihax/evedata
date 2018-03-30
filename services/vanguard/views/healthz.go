package views

import (
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
)

func init() {
	vanguard.AddRoute("healthz", "GET", "/healthz", healthz)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	g := vanguard.GlobalsFromContext(r.Context())

	// Check DB connection
	if err := g.Db.Ping(); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		return
	}

	// Check redis connection
	conn := g.Cache.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		return
	}
}
