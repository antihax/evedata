package views

import (
	"net/http"
	"runtime/pprof"

	"github.com/antihax/evedata/appContext"

	"github.com/antihax/evedata/server"
)

func init() {
	evedata.AddRoute("debug", "GET", "/debug/debugGoRoutines", debugGoRoutines)

}

func debugGoRoutines(c *appContext.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	setCache(w, 0)
	pprof.Lookup("goroutine").WriteTo(w, 1)
	return http.StatusOK, nil
}
