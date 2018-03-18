package tokenstore

import (
	"testing"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/goesi"
	"golang.org/x/oauth2"
)

func TestTokenStore(t *testing.T) {
	sql := sqlhelper.NewTestDatabase()
	models.SetDatabase(sql)
	err := models.AddCRESTToken(133, 133, "dude", &oauth2.Token{
		AccessToken:  "FAKE",
		RefreshToken: "So Fake",
		Expiry:       time.Now().Add(time.Hour * 100000),
		TokenType:    "Bearer"},
		"", "ownerhash", 0, 0, 0)
	if err != nil {
		t.Fatal(err)
		return
	}

	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()
	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(redis)

	// Setup an authenticator for our private token
	auth := goesi.NewSSOAuthenticator(cache, "fake", "reeeeely fake", "",
		[]string{"esi-universe.read_structures.v1",
			"esi-search.search_structures.v1",
			"esi-markets.structure_markets.v1"})

	ts := NewTokenStore(redis, sql, auth)

	err = ts.SetToken(133, 133, &oauth2.Token{
		RefreshToken: "fake",
		AccessToken:  "really fake",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}

	token, err := ts.GetToken(133, 133)
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "really fake" {
		t.Fatal("Token is incorrect 1")
	}

	err = ts.SetToken(133, 133, &oauth2.Token{
		RefreshToken: "fake",
		AccessToken:  "really very fake",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}

	tokensource, err := ts.GetTokenSource(133, 133)
	if err != nil {
		t.Fatal(err)
	}
	tok, err := tokensource.Token()
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "really very fake" {
		t.Fatal("Token is incorrect 2", tok.AccessToken)
	}
	err = ts.invalidateTokenCache(133, 133)
	if err != nil {
		t.Fatal(err)
	}
	token, err = ts.GetToken(133, 133)
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "really very fake" {
		t.Fatal("Token is incorrect 3")
	}
}
