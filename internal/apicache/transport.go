package apicache

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strconv"

	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// APICacheTransport to chain into the HTTPClient to gather statistics.
type APICacheTransport struct {
	Transport http.RoundTripper
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *APICacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Loop until success
	tries := 0
	for {
		// Tickup retry counter
		tries++

		// Time our response
		start := time.Now()

		// Run the request.
		res, httperr := t.Transport.RoundTrip(req)

		metricAPICalls.With(
			prometheus.Labels{"host": req.Host},
		).Observe(float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond))

		// We got a response
		if res != nil {
			// Get the ESI error information
			resetS := res.Header.Get("x-esi-error-limit-reset")
			tokensS := res.Header.Get("x-esi-error-limit-remain")

			// If we cannot decode this is likely from another source.
			esiRateLimiter := true
			reset, err := strconv.ParseFloat(resetS, 64)
			if err != nil {
				esiRateLimiter = false
			}

			// Tick up and log any errors
			if res.StatusCode >= 400 {
				metricAPIErrors.Inc()

				// Backoff
				sleep := 60 * time.Second
				if res.StatusCode == 401 { // Something went really wrong
					sleep = 120 * time.Second
				} else if esiRateLimiter { // Sleep based on error rate.
					sleep = time.Duration(reset+(rand.Float64()*15.0)) * time.Second
				} else if !esiRateLimiter { // Not an ESI error
					sleep = time.Second * time.Duration(tries) * 5
				}
				if sleep < time.Second {
					sleep = time.Second + (time.Duration(rand.Float32()) * time.Second)
				}
				if sleep > time.Second*60 {
					sleep = time.Second * 60
				}

				log.Printf("Try: %d Sleep: %d St: %d Res: %s Tok: %s - %s\n", tries, sleep/time.Second, res.StatusCode, resetS, tokensS, req.URL)

				// Dump data for important errors // !esiRateLimiter &&
				if res.StatusCode >= 400 {
					dump, _ := httputil.DumpResponse(res, true)
					fmt.Printf("%s\n\n", dump)
				}
				// Get out for "our bad" statuses
				if res.StatusCode >= 400 && res.StatusCode < 420 || res.StatusCode == 422 {
					if res.StatusCode != 403 {
						dump, _ := httputil.DumpRequest(req, true)
						fmt.Printf("%s\n\n", dump)
						log.Printf("Giving up %d %s\n", res.StatusCode, req.URL)
					}
					return res, httperr
				}

				time.Sleep(sleep)
			}

			if tries > 10 && res.StatusCode >= 400 {
				log.Printf("Too many tries %d %s\n", res.StatusCode, req.URL)
				return res, err
			}

		} else {
			return res, httperr
		}
		if res.StatusCode >= 200 && res.StatusCode < 400 {
			return res, httperr
		}
	}
}

var (
	metricAPICalls = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "api",
		Name:      "calls",
		Help:      "API call statistics.",
		Buckets:   prometheus.ExponentialBuckets(10, 1.45, 20),
	},
		[]string{"host"},
	)

	metricAPIErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "evedata",
		Subsystem: "api",
		Name:      "errors",
		Help:      "Count of API errors.",
	})
)

func init() {
	prometheus.MustRegister(
		metricAPICalls,
		metricAPIErrors,
	)
}
