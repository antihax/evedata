package apicache

import (
	"net/http"
)

var apiTransportLimiter chan bool

func init() {
	// concurrency limiter
	// 100 concurrent requests should fill 1 connection
	apiTransportLimiter = make(chan bool, 100)
}

// LimitedTransport limits concurrent requests to one connection.
type LimitedTransport struct {
	Transport http.RoundTripper
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *LimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Limit concurrency
	apiTransportLimiter <- true

	// Free the worker
	defer func() { <-apiTransportLimiter }()

	return t.Transport.RoundTrip(req)
}
