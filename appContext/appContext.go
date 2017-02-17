package appContext

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/antihax/evedata/config"
	"github.com/antihax/evedata/eveapi"
	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"
	"github.com/gregjones/httpcache"
	httpredis "github.com/gregjones/httpcache/redis"

	"golang.org/x/oauth2"
	gsr "gopkg.in/boj/redistore.v1"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// AppContext provides access to handles throughout the app.
type AppContext struct {
	Conf           *config.Config       // App Configuration
	Db             *sqlx.DB             // EVE Database
	Store          *gsr.RediStore       // Redis session store.
	EVE            *eveapi.EVEAPIClient // EVE API Client
	HTTPClient     *http.Client         // Redis Cached HTTP client
	Cache          *redis.Pool          // Redis connection Pool for HTTP Cache and session store.
	ESI            *goesi.APIClient
	ESIPublicToken oauth2.TokenSource

	// Since we need to combine data from multiple characters, we use
	// one authenticator for the site to act as the main authentication.
	// second will allow for many alt characters under the main.
	SSOAuthenticator          *eveapi.SSOAuthenticator // CREST authenticator for site authentication
	TokenAuthenticator        *eveapi.SSOAuthenticator // CREST authenticator for site functionality
	ESIBootstrapAuthenticator *eveapi.SSOAuthenticator // CREST authenticator for site functionality
}

func NewTestAppContext() AppContext {
	ctx := AppContext{}

	conf := config.Config{}
	ctx.Conf = &conf
	conf.EVEConsumer.Consumers = 10
	conf.EVEConsumer.ZKillEnabled = false

	database, err := models.SetupDatabase("mysql", "root@tcp(127.0.0.1:3306)/eve?allowOldPasswords=1&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	ctx.Db = database

	// Build the redis pool
	ctx.Cache = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	// Nuke anything in redis incase we have a flood of trash
	r := ctx.Cache.Get()
	r.Do("FLUSHALL")
	r.Close()

	// Create a Redis http client for the CCP APIs.
	transportCache := httpcache.NewTransport(httpredis.NewWithClient(ctx.Cache.Get()))

	// Attach a basic transport with our chained custom transport.
	transportCache.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment, MaxIdleConnsPerHost: 5}

	ctx.HTTPClient = &http.Client{Transport: transportCache}
	if ctx.HTTPClient == nil {
		log.Fatalln("client is null")
	}

	// Setup the EVE ESI Client
	ctx.ESI = goesi.NewAPIClient(ctx.HTTPClient, "EVEData.Org Test Client (If you can see me.. something broke)")
	ctx.ESI.ChangeBasePath("http://127.0.0.1:8080/latest")

	// Create a memcached session store.
	ctx.Store, err = gsr.NewRediStoreWithPool(ctx.Cache, []byte("SOME FAKE RANDOM KEY"))
	if err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	// Register structs for storage
	gob.Register(oauth2.Token{})
	gob.Register(eveapi.CRESTToken{})
	gob.Register(eveapi.VerifyResponse{})

	// Anonymous EVE API & Crest Client
	ctx.EVE = eveapi.NewEVEAPIClient(ctx.HTTPClient)

	// Setup the Token authenticator, this handles sub characters.
	tokenScopes := []string{
		eveapi.ScopeCharacterContractsRead,
		eveapi.ScopeCharacterMarketOrdersRead,
		eveapi.ScopeCharacterResearchRead,
		eveapi.ScopeCharacterWalletRead,
		"esi-assets.read_assets.v1",
		"esi-characters.read_contacts.v1",
		"esi-characters.write_contacts.v1",
		"esi-corporations.read_corporation_membership.v1",
		"esi-location.read_location.v1",
		"esi-location.read_ship_type.v1",
		"esi-planets.manage_planets.v1",
		"esi-search.search_structures.v1",
		"esi-skills.read_skills.v1",
		"esi-ui.open_window.v1",
		"esi-ui.write_waypoint.v1",
		"esi-universe.read_structures.v1",
		"esi-wallet.read_character_wallet.v1",
	}

	// take care to never actually make real requests on this.
	ctx.TokenAuthenticator = eveapi.NewSSOAuthenticator(
		ctx.HTTPClient,
		"123545",
		"PLEASE IGNORE",
		"I DO NOTHING",
		tokenScopes)

	return ctx
}
