package evedata

import (
	"encoding/gob"
	"evedata/appContext"
	"evedata/config"
	"evedata/discord"
	"evedata/emdrConsumer"
	"evedata/esi"
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

// Load the server
func GoServer() {
	var err error

	// Make a new app context.
	ctx := &appContext.AppContext{}

	// Read configuation.
	if ctx.Conf, err = config.ReadConfig(); err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Build the redis pool
	ctx.Cache = redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ctx.Conf.Redis.Address)
			if err != nil {
				return nil, err
			}
			if ctx.Conf.Redis.Password != "" {
				if _, err := c.Do("AUTH", ctx.Conf.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}

	// Build Connection Pool
	if ctx.Db, err = models.SetupDatabase(ctx.Conf.Database.Driver, ctx.Conf.Database.Spec); err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	// Setup the SSO authenticator, this is the main login.
	ssoScopes := []string{
		eveapi.ScopeCharacterKillsRead, // Temporary
		eveapi.ScopeCharacterLocationRead,
		eveapi.ScopeCharacterNavigationWrite,
		eveapi.ScopeRemoteClientUI,
	}

	ctx.SSOAuthenticator = eveapi.NewSSOAuthenticator(ctx.Conf.CREST.SSO.ClientID,
		ctx.Conf.CREST.SSO.SecretKey,
		ctx.Conf.CREST.SSO.RedirectURL,
		ssoScopes)

	// Setup the Token authenticator, this handles sub characters.
	tokenScopes := []string{
		eveapi.ScopeCharacterContactsRead,
		eveapi.ScopeCharacterContactsWrite,
	}

	ctx.TokenAuthenticator = eveapi.NewSSOAuthenticator(ctx.Conf.CREST.Token.ClientID,
		ctx.Conf.CREST.Token.SecretKey,
		ctx.Conf.CREST.Token.RedirectURL,
		tokenScopes)

	// Create a Redis http client for the CCP APIs.
	ctx.TransportCache = httpcache.NewTransport(httpredis.NewWithClient(ctx.Cache.Get()))
	ctx.TransportCache.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment, MaxIdleConnsPerHost: 5}

	ctx.HTTPClient = &http.Client{Transport: ctx.TransportCache}

	ctx.ESI = esi.NewAPIClient(ctx.HTTPClient)

	ctx.ESIBootstrapAuthenticator = eveapi.NewSSOAuthenticator(ctx.Conf.CREST.ESIAccessToken.ClientID,
		ctx.Conf.CREST.ESIAccessToken.SecretKey,
		ctx.Conf.CREST.ESIAccessToken.RedirectURL,
		[]string{"esi-universe.read_structures.v1",
			"esi-search.search_structures.v1"})

	token := &eveapi.CRESTToken{
		AccessToken:  ctx.Conf.CREST.ESIAccessToken.AccessToken,
		TokenType:    ctx.Conf.CREST.ESIAccessToken.TokenType,
		RefreshToken: ctx.Conf.CREST.ESIAccessToken.RefreshToken,
		Expiry:       ctx.Conf.CREST.ESIAccessToken.Expiry,
	}
	ctx.ESIPublicToken, err = ctx.ESIBootstrapAuthenticator.TokenSource(ctx.HTTPClient, token)

	if err != nil {
		log.Fatalf("Error starting bootstrap ESI client: %v", err)
	}

	// Create a memcached session store.
	ctx.Store, err = gsr.NewRediStoreWithPool(&ctx.Cache, []byte(ctx.Conf.Store.Key))
	if err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	ctx.Store.Options.Domain = ctx.Conf.Store.Domain

	// Register structs for storage
	gob.Register(oauth2.Token{})
	gob.Register(eveapi.CRESTToken{})
	gob.Register(eveapi.VerifyResponse{})

	// Anonymous EVE API & Crest Client
	ctx.EVE = eveapi.NewAnonymousClient(ctx.HTTPClient)

	// Set our logging flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if ctx.Conf.EMDRCrestBridge.Enabled {
		log.Println("Starting EMDR <- CREST Bridge")
		go emdrConsumer.GoEMDRCrestBridge(ctx)
	}

	if ctx.Conf.EVEConsumer.Enabled {
		log.Println("Starting EVE Consumer")
		eC := eveConsumer.NewEVEConsumer(ctx)
		eC.RunConsumer()
		defer eC.StopConsumer()
	}

	if ctx.Conf.Discord.Enabled {
		go discord.GoDiscordBot(ctx)
	}

	// Allocate the routes
	rtr := NewRouter(ctx)

	log.Printf("EveData Listening port 3000...\n")
	http.ListenAndServe(":3000", context.ClearHandler(rtr))
	log.Printf("EveData Quitting..\n")
}
