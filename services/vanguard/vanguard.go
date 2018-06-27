package vanguard

import (
	"encoding/gob"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/discordauth"
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
	stop        chan bool
	wg          *sync.WaitGroup
	OutQueue    *redisqueue.RedisQueue
	ESI         *goesi.APIClient
	Db          *sqlx.DB
	Cache       *redis.Pool
	HTTPClient  *http.Client
	Store       *gsr.RediStore
	TokenStore  *tokenstore.TokenStore
	Conservator *rpc.Client

	// authentication
	token                *oauth2.TokenSource
	TokenAuthenticator   *goesi.SSOAuthenticator
	SSOAuthenticator     *goesi.SSOAuthenticator
	DiscordAuthenticator *discordauth.Authenticator
}

var globalVanguard *Vanguard

func NewVanguard(redis *redis.Pool, ledis *redis.Pool, db *sqlx.DB) *Vanguard {
	// Don't allow more than one to be created
	if globalVanguard != nil {
		return globalVanguard
	}
	models.SetDatabase(db)
	// Register structs for storage
	gob.Register(oauth2.Token{})
	gob.Register(goesi.VerifyResponse{})

	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(ledis)

	// Create our ESI API Client
	esi := goesi.NewAPIClient(cache, "EVEData-API-Vanguard")

	// Setup an authenticator for our user tokens
	tauth := goesi.NewSSOAuthenticator(cache, os.Getenv("ESI_CLIENTID_TOKENSTORE"), os.Getenv("ESI_SECRET_TOKENSTORE"), "https://"+os.Getenv("DOMAIN")+"/X/eveTokenAnswer", []string{})

	// Setup an authenticator for our SSO token
	ssoauth := goesi.NewSSOAuthenticator(cache, os.Getenv("ESI_CLIENTID_SSO"), os.Getenv("ESI_SECRET_SSO"), "https://"+os.Getenv("DOMAIN")+"/X/eveSSOAnswer", []string{})

	// Setup an authenticator for Discord
	dauth := discordauth.NewAuthenticator(cache, os.Getenv("DISCORD_CLIENTID"), os.Getenv("DISCORD_SECRET"), "https://"+os.Getenv("DOMAIN")+"/X/discordAnswer", []string{"identify", "guilds.join"})

	// Build our private token
	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: os.Getenv("ESI_REFRESHKEY"),
		TokenType:    "Bearer",
	}
	token, err := tauth.TokenSource(tok)
	if err != nil {
		log.Fatalln(err)
	}

	// create Token Store
	tokenStore := tokenstore.NewTokenStore(redis, db, tauth)

	// Create a redis session store.
	store, err := gsr.NewRediStoreWithPool(redis, []byte(os.Getenv("COOKIE_SECRET")))
	if err != nil {
		log.Fatalf("Cannot build redis store: %v", err)
	}
	// Set options for the store
	store.SetMaxLength(1024 * 100)
	store.Options.Domain = "evedata.org"

	// Setup a new Vanguard
	globalVanguard = &Vanguard{
		stop: make(chan bool),
		wg:   &sync.WaitGroup{},
		OutQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),

		HTTPClient:           cache,
		SSOAuthenticator:     ssoauth,
		TokenAuthenticator:   tauth,
		DiscordAuthenticator: dauth,
		ESI:                  esi,
		Store:                store,
		Db:                   db,
		Cache:                redis,

		token:      &token,
		TokenStore: tokenStore,
	}

	if err := globalVanguard.RPCConnect(); err != nil {
		log.Fatalln(err)
	}

	return globalVanguard
}

// Close the  service
func (s *Vanguard) Close() {
	close(s.stop)
	s.wg.Wait()
}

// RPCall calls remote procedures
func (s *Vanguard) RPCall(method string, in interface{}, out interface{}) error {
	for {
		err := s.Conservator.Call(method, in, out)
		if err == rpc.ErrShutdown {
			err := s.RPCConnect()
			if err != nil {
				log.Printf("lost rpc connection: %s", err)
				time.Sleep(time.Millisecond * 50)
			}
			continue
		}
		return err
	}
}

// Close the service
func (s *Vanguard) RPCConnect() error {
	var err error
	s.Conservator, err = rpc.DialHTTP("tcp", "conservator.evedata:3001")
	if err != nil {
		s.Conservator, err = rpc.DialHTTP("tcp", "conservator.evedata:32003")
		if err != nil {
			return err
		}
	}

	return nil
}
