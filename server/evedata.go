package evedata

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/config"
	"github.com/antihax/evedata/discord"
	"github.com/antihax/evedata/esi"
	"github.com/antihax/evedata/eveConsumer"
	"github.com/antihax/evedata/eveapi"
	"github.com/antihax/evedata/models"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/context"
	"github.com/gregjones/httpcache"
	httpredis "github.com/gregjones/httpcache/redis"
	"golang.org/x/oauth2"
	gsr "gopkg.in/boj/redistore.v1"
)

var ctx appContext.AppContext

func GetContext() *appContext.AppContext {
	return &ctx
}

// Load the server
func GoServer() {
	var err error

	// Make a new app context.

	// Read configuation.
	if ctx.Conf, err = config.ReadConfig(); err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Build the redis pool
	ctx.Cache = &redis.Pool{
		MaxIdle:     10,
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
		"esi-assets.read_assets.v1",
		"esi-bookmarks.read_character_bookmarks.v1",
		"esi-calendar.read_calendar_events.v1",
		"esi-calendar.respond_calendar_events.v1",
		"esi-characters.read_contacts.v1",
		"esi-characters.write_contacts.v1",
		"esi-clones.read_clones.v1",
		"esi-corporations.read_corporation_membership.v1",
		"esi-fleets.read_fleet.v1",
		"esi-fleets.write_fleet.v1",
		"esi-killmails.read_killmails.v1",
		"esi-location.read_location.v1",
		"esi-location.read_ship_type.v1",
		"esi-planets.manage_planets.v1",
		"esi-search.search_structures.v1",
		"esi-skills.read_skillqueue.v1",
		"esi-skills.read_skills.v1",
		"esi-ui.open_window.v1",
		"esi-ui.write_waypoint.v1",
		"esi-universe.read_structures.v1",
		"esi-wallet.read_character_wallet.v1",
	}

	ctx.TokenAuthenticator = eveapi.NewSSOAuthenticator(ctx.Conf.CREST.Token.ClientID,
		ctx.Conf.CREST.Token.SecretKey,
		ctx.Conf.CREST.Token.RedirectURL,
		tokenScopes)

	// Create a Redis http client for the CCP APIs.
	ctx.TransportCache = httpcache.NewTransport(httpredis.NewWithClient(ctx.Cache.Get()))
	ctx.TransportCache.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment, MaxIdleConnsPerHost: 5}

	ctx.HTTPClient = &http.Client{Transport: ctx.TransportCache}

	ctx.ESI = esi.NewAPIClient(ctx.HTTPClient, ctx.Conf.UserAgent)

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
	ctx.Store, err = gsr.NewRediStoreWithPool(ctx.Cache, []byte(ctx.Conf.Store.Key))
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

	if ctx.Conf.Discord.Enabled {
		go discord.GoDiscordBot(&ctx)
	}

	if ctx.Conf.EVEConsumer.Enabled {
		log.Println("Starting EVE Consumer")
		eC := eveConsumer.NewEVEConsumer(&ctx)
		eC.RunConsumer()
		defer eC.StopConsumer()
	}
	// Allocate the routes
	rtr := NewRouter(&ctx)

	log.Printf("EveData Listening port 3000...\n")
	http.ListenAndServe(":3000", context.ClearHandler(rtr))
	log.Printf("EveData Quitting..\n")
}
