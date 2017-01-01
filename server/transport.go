package evedata

import (
	"crypto/rand"
	"net/http"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/models"
)

// Custom transport to chain into the HTTPClient to gather statistics.
type transport struct {
	next *http.Transport
	ctx  *appContext.AppContext
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to keep track
// of the current request.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	redisCon := ctx.Cache.Get()
	defer redisCon.Close()

	// Make a random hash to store the time to redis
	b := make([]byte, 32)
	rand.Read(b)

	redisCon.Do("ZADD", "EVEDATA_HTTPCount", time.Now().UTC().Unix(), b)
	res, err := t.next.RoundTrip(req)
	if res.StatusCode >= 400 {
		redisCon.Do("ZADD", "EVEDATA_HTTPErrorCount", time.Now().UTC().Unix(), b)
		models.AddHTTPError(req, res)
	}
	return res, err
}
