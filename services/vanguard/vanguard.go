package vanguard

import (
	"encoding/gob"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/goesi"
	gsr "github.com/antihax/redistore"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

type Vanguard struct {
	stop       chan bool
	wg         *sync.WaitGroup
	OutQueue   *redisqueue.RedisQueue
	ESI        *goesi.APIClient
	Db         *sqlx.DB
	Cache      *redis.Pool
	HTTPClient *http.Client
	Store      *gsr.RediStore
	TokenStore *tokenstore.TokenStore

	// authentication
	token                *oauth2.TokenSource
	PrivateAuthenticator *goesi.SSOAuthenticator
	TokenAuthenticator   *goesi.SSOAuthenticator
	SSOAuthenticator     *goesi.SSOAuthenticator
}

var globalVanguard *Vanguard

func NewVanguard(redis *redis.Pool, db *sqlx.DB, clientID, secret, refresh, tokenClientID, tokenSecret string, ssoClientID, ssoSecret string, storeKey string) *Vanguard {
	// Don't allow more than one to be created
	if globalVanguard != nil {
		return globalVanguard
	}
	models.SetDatabase(db)
	// Register structs for storage
	gob.Register(oauth2.Token{})
	gob.Register(goesi.VerifyResponse{})

	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(redis)

	// Create our ESI API Client
	esi := goesi.NewAPIClient(cache, "EVEData-API-Vanguard")

	// Setup an authenticator for our user tokens
	tauth := goesi.NewSSOAuthenticator(cache, tokenClientID, tokenSecret, "", []string{})

	// Setup an authenticator for our private token
	pauth := goesi.NewSSOAuthenticator(cache, clientID, secret, "",
		[]string{"esi-universe.read_structures.v1",
			"esi-search.search_structures.v1",
			"esi-markets.structure_markets.v1"})

	// Setup an authenticator for our SSO token
	ssoauth := goesi.NewSSOAuthenticator(cache, ssoClientID, ssoSecret, "", []string{})

	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}

	// Build our private token
	token, err := pauth.TokenSource(tok)
	if err != nil {
		log.Fatalln(err)
	}

	tokenStore := tokenstore.NewTokenStore(redis, db, tauth)

	// Create a redis session store.
	store, err := gsr.NewRediStoreWithPool(redis, []byte(storeKey))
	if err != nil {
		log.Fatalf("Cannot build redis store: %v", err)
	}

	// Setup a new Vanguard
	globalVanguard = &Vanguard{
		stop: make(chan bool),
		wg:   &sync.WaitGroup{},
		OutQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),

		HTTPClient:           cache,
		PrivateAuthenticator: pauth,
		SSOAuthenticator:     ssoauth,
		TokenAuthenticator:   tauth,
		ESI:                  esi,
		Store:                store,
		Db:                   db,
		Cache:                redis,
		token:                &token,
		TokenStore:           tokenStore,
	}

	return globalVanguard
}

// Close the hammer service
func (s *Vanguard) Close() {
	close(s.stop)
	s.wg.Wait()
}
