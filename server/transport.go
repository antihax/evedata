package evedata

import (
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/antihax/evedata/appContext"
)

// Custom transport to chain into the HTTPClient to gather statistics.
type transport struct {
	next      *http.Transport
	ctx       *appContext.AppContext
	errorRate uint32
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	redisCon := ctx.Cache.Get()
	defer redisCon.Close()

	// Make a random hash to store the time to redis
	b := make([]byte, 32)
	rand.Read(b)

	redisCon.Do("ZADD", "EVEDATA_HTTPCount", time.Now().UTC().Unix(), b)

	// Do the request.
	res, err := t.next.RoundTrip(req)

	// We got a non-recoverable error.
	if res != nil {
		if res.StatusCode >= 400 {
			redisCon.Do("ZADD", "EVEDATA_HTTPErrorCount", time.Now().UTC().Unix(), b)

			// Tick up the error rate and sleep proportionally to the error count.
			if res.StatusCode >= 500 {
				if t.errorRate < 60 {
					atomic.AddUint32(&t.errorRate, 1)
				}
				time.Sleep(time.Second * time.Duration(t.errorRate))
			}
		} else {
			// Tick down the error rate.
			if t.errorRate > 0 {
				atomic.AddUint32(&t.errorRate, ^uint32(0))
			}
		}
	}

	return res, err
}
