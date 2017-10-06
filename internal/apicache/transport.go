package apicache

import (
	"net/http"
	"strconv"

	"time"
)

// Custom transport to chain into the HTTPClient to gather statistics.
type transport struct {
	next *http.Transport
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Loop until success
	tries := 0
	for {
		tries++
		esiRateLimiter := true
		// Do the request.
		res, err := t.next.RoundTrip(req)

		if res != nil {
			// Get the ESI error information
			resetS := res.Header.Get("x-esi-error-limit-reset")
			tokensS := res.Header.Get("x-esi-error-limit-remain")

			// If we cannot decode this is likely from another source.
			reset, err := strconv.Atoi(resetS)
			if err != nil {
				esiRateLimiter = false
			}
			tokens, err := strconv.Atoi(tokensS)
			if err != nil {
				esiRateLimiter = false
			}

			// Sleep to prevent hammering CCP ESI if there are excessive errors
			if esiRateLimiter {
				time.Sleep(time.Duration(float64(reset*3) * (1 - (float64(tokens) / 100))))
			}

			if res.StatusCode == 420 {
				time.Sleep(time.Duration(reset))
			}

			if res.StatusCode == 420 || res.StatusCode >= 500 || res.StatusCode == 0 {
				// break out after 10 tries
				if tries > 10 {
					return res, err
				}
				if !esiRateLimiter {
					time.Sleep(time.Duration(tries * 2))
				}
				continue
			}
		}

		return res, err
	}
}
