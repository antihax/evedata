package evedata

import (
	"net/http"
	"sync/atomic"

	"time"

	"github.com/ScaleFT/monotime"
	"github.com/antihax/evedata/appContext"
	"github.com/prometheus/client_golang/prometheus"
)

// Custom transport to chain into the HTTPClient to gather statistics.
type transport struct {
	next *http.Transport
	ctx  *appContext.AppContext
}

var errorRate int32

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := monotime.Now()

	// Do the request.
	res, err := t.next.RoundTrip(req)

	duration := monotime.Duration(start, monotime.Now())
	metricApiCalls.With(
		prometheus.Labels{"host": req.Host},
	).Observe(float64(duration / time.Millisecond))

	// We got a non-recoverable error.
	if res != nil {

		if res.StatusCode >= 400 || res.StatusCode == 0 {
			//	models.AddHTTPError(req, res)
			metricApiErrors.Inc()
			// Tick up the error rate and sleep proportionally to the error count.
			if res.StatusCode >= 500 || res.StatusCode == 000 {
				errors := atomic.LoadInt32(&errorRate)
				if errors < 60 {
					atomic.AddInt32(&errorRate, 1)
					errors++
				} else if errors > 60 {
					atomic.StoreInt32(&errorRate, 60)
				}
				time.Sleep(time.Second * time.Duration(errors))
			}
		} else {
			// Tick down the error rate.
			errors := atomic.LoadInt32(&errorRate)
			if errors > 0 {
				atomic.AddInt32(&errorRate, ^int32(0))
			} else if errors < 0 {
				atomic.StoreInt32(&errorRate, 0)
			}
		}
	}

	return res, err
}

var (
	metricApiCalls = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "api",
		Name:      "calls",
		Help:      "API call statistics.",
		Buckets:   prometheus.LinearBuckets(0, 50, 25),
	},
		[]string{"host"},
	)

	metricApiErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "evedata",
		Subsystem: "api",
		Name:      "errors",
		Help:      "Count of API errors.",
	})
)

func init() {
	prometheus.MustRegister(
		metricApiCalls,
		metricApiErrors,
	)
}
