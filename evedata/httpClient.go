package evedata

import (
	"net"
	"net/http"
	"time"

	"log"

	"github.com/antihax/httpcache"
	httpredis "github.com/antihax/httpcache/redis"
	"github.com/garyburd/redigo/redis"
)

func setupHTTPClient(cache *redis.Pool) *http.Client {
	// Create a Redis http client for the CCP APIs.
	transportCache := httpcache.NewTransport(httpredis.NewWithClient(cache))

	// Attach a basic transport with our chained custom transport.
	transportCache.Transport = &transport{&http.Transport{
		MaxIdleConns: 200,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 0,
		MaxIdleConnsPerHost:   20,
	}, &ctx}

	client := &http.Client{Transport: transportCache}
	if client == nil {
		log.Fatalf("http client is null")
		panic("http client is null")
	}
	return client
}
