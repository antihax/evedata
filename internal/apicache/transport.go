package apicache

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Custom transport to chain into the HTTPClient to gather statistics.
type ApiCacheTransport struct {
	next *http.Transport
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *ApiCacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Loop until success
	tries := 0
	for {
		// Tickup retry counter
		tries++

		// Time our response
		start := time.Now()

		// Run the request.
		res, err := t.next.RoundTrip(req)

		metricAPICalls.With(
			prometheus.Labels{"host": req.Host},
		).Observe(float64(time.Since(start).Nanoseconds()) / float64(time.Millisecond))

		// We got a response
		if res != nil {
			// Get the ESI error information
			resetS := res.Header.Get("x-esi-error-limit-reset")
			tokensS := res.Header.Get("x-esi-error-limit-remain")

			// Tick up and log any errors
			if res.StatusCode >= 400 {
				metricAPIErrors.Inc()
				log.Printf("St: %d Res: %s Tok: %s - %s\n", res.StatusCode, resetS, tokensS, req.URL)
			}

			// If we cannot decode this is likely from another source.
			esiRateLimiter := true
			reset, err := strconv.ParseFloat(resetS, 64)
			if err != nil {
				esiRateLimiter = false
			}
			tokens, err := strconv.ParseFloat(tokensS, 64)
			if err != nil {
				esiRateLimiter = false
			}

			// Backoff
			if res.StatusCode == 420 { // Something went wrong
				duration := reset * ((1 + rand.Float64()) * 5)
				time.Sleep(time.Duration(duration) * time.Second)
			} else if esiRateLimiter { // Sleep based on error rate.
				percentRemain := 1 - (tokens / 100)
				duration := reset * percentRemain * (1 + rand.Float64())
				time.Sleep(time.Second * time.Duration(duration))
			} else if !esiRateLimiter { // Not an ESI error
				time.Sleep(time.Second * time.Duration(tries))
			}

			// Get out for "our bad" statuses
			if res.StatusCode >= 400 && res.StatusCode < 420 {
				if res.StatusCode != 403 {
					log.Printf("Giving up %d %s\n", res.StatusCode, req.URL)
				}
				return res, err
			}

			if tries > 10 {
				log.Printf("Too many tries %d %s\n", res.StatusCode, req.URL)
				return res, err
			}
		} else {
			return res, err
		}
		if res.StatusCode >= 200 && res.StatusCode < 400 {
			return res, err
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
