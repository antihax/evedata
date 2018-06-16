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

		esiRateLimiter := true

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

			// Early out for "our bad" statuses
			if (res.StatusCode >= 400 && res.StatusCode < 420) || res.StatusCode == 500 {
				return res, err
			}

			// If we cannot decode this is likely from another source.
			reset, err := strconv.Atoi(resetS)
			if err != nil {
				esiRateLimiter = false
			}
			tokens, err := strconv.Atoi(tokensS)
			if err != nil {
				esiRateLimiter = false
			}

			// Backoff
			if res.StatusCode == 420 { // Something went wrong
				time.Sleep(time.Duration(reset)*time.Second + time.Duration(rand.Intn(int(time.Second*5))))
			} else if res.StatusCode == 429 { // SNAFU
				time.Sleep(time.Second * time.Duration(60+rand.Intn(30)))
			} else if esiRateLimiter { // Sleep based on error rate.
				time.Sleep(time.Second * time.Duration(float64(reset)*2*(1-(float64(tokens)/100))))
			} else if !esiRateLimiter { // Not an ESI error
				time.Sleep(time.Second * time.Duration(tries))
			}

			// break out after 10 tries
			if res.StatusCode == 420 || res.StatusCode == 429 || res.StatusCode >= 500 || res.StatusCode == 0 {
				if tries > 10 {
					return res, err
				}
			}
		}
		return res, err
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
