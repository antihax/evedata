// Package tailor handles fitting attributes to the database
package tailor

import (
	"io/ioutil"
	"log"
	"net/http"
)

// ApiTransport custom transport to chain into the HTTPClient to gather statistics.
type ApiTransport struct {
	Next *http.Transport
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *ApiTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Run the request and time the response
	res, triperr := t.Next.RoundTrip(req)

	// We got a response
	if res != nil {
		// Tick up and log any errors
		if res.StatusCode >= 500 {
			b, err := ioutil.ReadAll(res.Body)
			log.Printf("St: %d %+v %s \n", res.StatusCode, string(b), err)
		}
	}

	return res, triperr
}
