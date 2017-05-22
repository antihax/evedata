package apicache

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/stretchr/testify/assert"
)

func TestAPICache(t *testing.T) {
	redis := redigohelper.ConnectRedisTestPool()

	// Setup a simple http server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		size, _ := strconv.Atoi(r.FormValue("size"))
		w.Header().Set("Cache-Control", "max-age=2592000")
		w.Write(make([]byte, size))
	}))
	defer server.Close()

	// Test the client
	client := CreateHTTPClientCache(redis)
	res, err := client.Get("http://" + server.Listener.Addr().String() + "/?size=30")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "", res.Header.Get("x-from-cache"))

	res, err = client.Get("http://" + server.Listener.Addr().String() + "/?size=30")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "1", res.Header.Get("x-from-cache"))
}
