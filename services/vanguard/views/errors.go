package views

import (
	"log"
	"net/http"

	"github.com/antihax/evedata/services/vanguard"
)

func init() {
	vanguard.AddNotFoundHandler(notFoundPage)
}

func httpErrCode(w http.ResponseWriter, err error, code int) {
	if err != nil {
		log.Printf("http error %s", err)
	}
	http.Error(w, http.StatusText(code), code)
}

func httpErr(w http.ResponseWriter, err error) {
	log.Printf("http error %s", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
