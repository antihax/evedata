package evedata

import (
	"encoding/gob"
	"evedata/appContext"
	"evedata/config"
	"evedata/eveConsumer"
	"evedata/eveapi"
	"evedata/models"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/context"
	"github.com/gregjones/httpcache"
	httpredis "github.com/gregjones/httpcache/redis"
	"golang.org/x/oauth2"
	gsr "gopkg.in/boj/redistore.v1"
)

func GoServer() {
	var err error

	// Make a new app context.
	ctx := &appContext.AppContext{}

	// Read configuation.
	if ctx.Conf, err = config.ReadConfig(); err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Build Connection Pool
	if ctx.Db, err = models.SetupDatabase(ctx.Conf.Database.Driver, ctx.Conf.Database.Spec); err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	scopes := []string{eveapi.ScopeCharacterContactsRead,
		eveapi.ScopeCharacterContactsWrite}

	ctx.SSOAuthenticator = eveapi.NewSSOAuthenticator(ctx.Conf.CREST.SSO.ClientID,
		ctx.Conf.CREST.SSO.SecretKey,
		ctx.Conf.CREST.SSO.RedirectURL,
		nil)

	ctx.TokenAuthenticator = eveapi.NewSSOAuthenticator(ctx.Conf.CREST.Token.ClientID,
		ctx.Conf.CREST.Token.SecretKey,
		ctx.Conf.CREST.Token.RedirectURL,
		scopes)

	ctx.Cache = redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ctx.Conf.Redis.Address)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", ctx.Conf.Redis.Password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}

	// Create a memcached http client for the CCP APIs.
	transport := httpcache.NewTransport(httpredis.NewWithClient(ctx.Cache.Get()))
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment, MaxIdleConnsPerHost: 5}

	ctx.HTTPClient = &http.Client{Transport: transport}

	// Create a memcached session store.
	ctx.Store, err = gsr.NewRediStoreWithPool(&ctx.Cache, []byte(ctx.Conf.Store.Key))
	if err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	ctx.Store.Options.Domain = ctx.Conf.Domain

	// Register structs for storage
	gob.Register(oauth2.Token{})
	gob.Register(eveapi.CRESTToken{})
	gob.Register(eveapi.VerifyResponse{})

	// Anonymous EVE API & Crest Client
	ctx.EVE = eveapi.NewAnonymousClient(ctx.HTTPClient)

	// Set our logging flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if ctx.Conf.EMDRCrestBridge.Enabled {
		log.Println("Starting EMDR <- Crest Bridge")
		go goEMDRCrestBridge(ctx)
	}

	eC := eveConsumer.NewEVEConsumer(ctx)
	eC.RunConsumer()
	defer eC.StopConsumer()

	// Allocate the routes
	rtr := NewRouter(ctx)

	log.Printf("EveData Listening port 3000...\n")
	http.ListenAndServe(":3000", context.ClearHandler(rtr))
}
