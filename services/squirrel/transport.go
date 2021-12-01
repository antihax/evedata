package squirrel

import (
	"net/http"
)

// LimiterTransport provides concurrency limiting to outbound http calls
type LimiterTransport struct {
	next http.RoundTripper
}

var semLimiterTransport chan bool

func init() {
	// concurrency limiter
	// 100 concurrent requests should fill 1 connection
	semLimiterTransport = make(chan bool, 100)
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to limit concurrency.
func (t LimiterTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Run the request.
	// Limit concurrency
	semLimiterTransport <- true
	//log.Printf("%+v\n", req.URL)

	// Free the worker
	defer func() { <-semLimiterTransport }()
	return t.next.RoundTrip(req)
}
