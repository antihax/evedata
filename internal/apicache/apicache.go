package apicache

import (
	"net"
	"net/http"
	"time"

	"github.com/antihax/httpcache"
	httpredis "github.com/antihax/httpcache/redis"
	"github.com/garyburd/redigo/redis"
)

// CreateHTTPClientCache creates an error limiting client with auto retry and redis cache
func CreateHTTPClientCache(redis *redis.Pool) *http.Client {

	// Create a Redis http client for the CCP APIs.
	transportCache := httpcache.NewTransport(httpredis.NewWithClient(redis))

	// Attach a basic transport with our chained custom transport.
	transportCache.Transport = &APICacheTransport{
		&http.Transport{
			MaxIdleConns: 100,
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 5 * 60 * time.Second,
				DualStack: true,
			}).DialContext,
			IdleConnTimeout:       5 * 60 * time.Second,
			TLSHandshakeTimeout:   60 * time.Second,
			ResponseHeaderTimeout: 60 * time.Second,
			ExpectContinueTimeout: 0,
			MaxIdleConnsPerHost:   3,
		},
	}

	client := &http.Client{Transport: transportCache}
	if client == nil {
		panic("http client is null")
	}
	return client
}

// CreateLimitedHTTPClientCache creates an error limiting client with auto retry, redis cache,
// and limit it to 100 connections
func CreateLimitedHTTPClientCache(redis *redis.Pool) *http.Client {

	// Create a Redis http client for the CCP APIs.
	transportCache := httpcache.NewTransport(httpredis.NewWithClient(redis))

	cache := &LimitedTransport{
		&APICacheTransport{
			&http.Transport{
				MaxIdleConns: 100,
				DialContext: (&net.Dialer{
					Timeout:   60 * time.Second,
					KeepAlive: 5 * 60 * time.Second,
					DualStack: true,
				}).DialContext,
				IdleConnTimeout:       5 * 60 * time.Second,
				TLSHandshakeTimeout:   60 * time.Second,
				ResponseHeaderTimeout: 60 * time.Second,
				ExpectContinueTimeout: 0,
				MaxIdleConnsPerHost:   3,
			},
		},
	}

	// Attach a basic transport with our chained custom transport.
	transportCache.Transport = cache

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
	t := &APICacheTransport{
		&http.Transport{
			MaxIdleConns: 100,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 5 * 60 * time.Second,
				DualStack: true,
			}).DialContext,
			IdleConnTimeout:       5 * 60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			ExpectContinueTimeout: 0,
			MaxIdleConnsPerHost:   3,
		},
	}

	client := &http.Client{Transport: t}
	if client == nil {
		panic("http client is null")
	}
	return client
}
