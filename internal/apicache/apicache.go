package apicache

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/httpcache"
	httpredis "github.com/antihax/httpcache/redis"
)

// CreateHTTPClientCache creates an error limiting client with auto retry and redis cache
func CreateHTTPClientCache() *http.Client {

	ledis := redigohelper.ConnectLedisProdPool()
	if ledis == nil {
		log.Fatal("ledis is nil")
	}
	// Create a Redis http client for the CCP APIs.
	transportCache := httpcache.NewTransport(httpredis.NewWithClient(ledis))

	// Attach a basic transport with our chained custom transport.
	transportCache.Transport = &ApiCacheTransport{
		&http.Transport{
			MaxIdleConns: 200,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 5 * 60 * time.Second,
				DualStack: true,
			}).DialContext,
			IdleConnTimeout:       5 * 60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			ExpectContinueTimeout: 0,
			MaxIdleConnsPerHost:   20,
		},
	}

	client := &http.Client{Transport: transportCache}
	if client == nil {
		panic("http client is null")
	}
	return client
}

// CreateHTTPClientCache creates an error limiting client with auto retry and no cache
func CreateHTTPClient() *http.Client {
	// Create a Redis http client for the CCP APIs.

	// Attach a basic transport with our chained custom transport.
	t := &ApiCacheTransport{
		&http.Transport{
			MaxIdleConns: 200,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 5 * 60 * time.Second,
				DualStack: true,
			}).DialContext,
			IdleConnTimeout:       5 * 60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			ExpectContinueTimeout: 0,
			MaxIdleConnsPerHost:   20,
		},
	}

	client := &http.Client{Transport: t}
	if client == nil {
		panic("http client is null")
	}
	return client
}
