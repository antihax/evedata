package tokenserver

import (
	"net/http"
	"testing"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/goesi"
)

func TestTokenServer(t *testing.T) {
	sql := sqlhelper.NewTestDatabase()

	redis := redigohelper.ConnectRedisTestPool()
	redConn := redis.Get()
	defer redConn.Close()
	redConn.Do("FLUSHALL")

	// Setup an authenticator
	auth := goesi.NewSSOAuthenticator(&http.Client{}, "FAAAK", "SOFAAAAAAAKE", "",
		[]string{"esi-universe.read_structures.v1",
			"esi-search.search_structures.v1",
			"esi-markets.structure_markets.v1"})
	auth.ChangeTokenURL("http://127.0.0.1:8080")
	auth.ChangeAuthURL("http://127.0.0.1:8080")

	ts := NewTokenServer(redis, sql, auth)
	go ts.Run()
}
