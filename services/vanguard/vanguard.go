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
	token              *oauth2.TokenSource
	TokenAuthenticator *goesi.SSOAuthenticator
	SSOAuthenticator   *goesi.SSOAuthenticator
}

var globalVanguard *Vanguard

func NewVanguard(redis *redis.Pool, db *sqlx.DB, refresh, tokenClientID, tokenSecret, ssoClientID, ssoSecret, storeKey, domain string) *Vanguard {
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
	tauth := goesi.NewSSOAuthenticator(cache, tokenClientID, tokenSecret, "https://"+domain+"/X/eveTokenAnswer", []string{})

	// Setup an authenticator for our SSO token
	ssoauth := goesi.NewSSOAuthenticator(cache, ssoClientID, ssoSecret, "https://"+domain+"/X/eveSSOAnswer", []string{})

	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}

	// Build our private token
	token, err := tauth.TokenSource(tok)
	if err != nil {
		log.Fatalln(err)
	}

	// create Token Store
	tokenStore := tokenstore.NewTokenStore(redis, db, tauth)

	// Create a redis session store.
	store, err := gsr.NewRediStoreWithPool(redis, []byte(storeKey))
	if err != nil {
		log.Fatalf("Cannot build redis store: %v", err)
	}
	// Set options for the store
	store.SetMaxLength(1024 * 10)
	store.Options.Domain = "evedata.org"

	// Setup a new Vanguard
	globalVanguard = &Vanguard{
		stop: make(chan bool),
		wg:   &sync.WaitGroup{},
		OutQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),

		HTTPClient:         cache,
		SSOAuthenticator:   ssoauth,
		TokenAuthenticator: tauth,
		ESI:                esi,
		Store:              store,
		Db:                 db,
		Cache:              redis,
		token:              &token,
		TokenStore:         tokenStore,
	}

	return globalVanguard
}

// Close the hammer service
func (s *Vanguard) Close() {
	close(s.stop)
	s.wg.Wait()
}
