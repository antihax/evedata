package evedata

import (
	"encoding/gob"
	"net"
	"os"

	"log"
	"net/http"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/config"
	"github.com/antihax/evedata/discord"
	"github.com/antihax/evedata/eveConsumer"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"

	"github.com/antihax/httpcache"
	httpredis "github.com/antihax/httpcache/redis"
	gsr "github.com/antihax/redistore"
	"github.com/gorilla/context"
	"golang.org/x/oauth2"
)

var ctx appContext.AppContext

// GetContext Returns the global appContext for EVEData Server
func GetContext() *appContext.AppContext {
	return &ctx
}

// GoServer Runs the EVEData Server
func GoServer() {
	var err error

	// Read configuation.
	if ctx.Conf, err = config.ReadConfig(); err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Build the redis pool
	ctx.Cache = setupRedis(GetContext())

	// Create a Redis http client for the CCP APIs.
	transportCache := httpcache.NewTransport(httpredis.NewWithClient(ctx.Cache))

	// Attach a basic transport with our chained custom transport.
	transportCache.Transport = &transport{&http.Transport{
		MaxIdleConns: 60,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		ExpectContinueTimeout: 0,
		MaxIdleConnsPerHost:   5,
	}, &ctx}

	// Build a HTTP Client pool this client will be shared with APIs for:
	//   - ESI
	//   - ZKillboard
	//   - EVE SSO
	//   - EVE CREST and XML

	ctx.HTTPClient = &http.Client{Transport: transportCache}
	if ctx.HTTPClient == nil {
		panic("http client is null")
	}

	// Build Connection Pool
	if ctx.Db, err = models.SetupDatabase(ctx.Conf.Database.Driver, ctx.Conf.Database.Spec); err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	/*r := ctx.Cache.Get()
	r.Do("FLUSHALL")
	r.Close()*/

	// Setup the SSO authenticator, this is the main login.
	ssoScopes := []string{}

	ctx.SSOAuthenticator = goesi.NewSSOAuthenticator(
		ctx.HTTPClient,
		ctx.Conf.CREST.SSO.ClientID,
		ctx.Conf.CREST.SSO.SecretKey,
		ctx.Conf.CREST.SSO.RedirectURL,
		ssoScopes)

	// Setup the Token authenticator, this handles sub characters.
	tokenScopes := models.GetCharacterScopes()
	ctx.TokenAuthenticator = goesi.NewSSOAuthenticator(
		ctx.HTTPClient,
		ctx.Conf.CREST.Token.ClientID,
		ctx.Conf.CREST.Token.SecretKey,
		ctx.Conf.CREST.Token.RedirectURL,
		tokenScopes)

	// Setup the EVE ESI Client
	ctx.ESI = goesi.NewAPIClient(ctx.HTTPClient, ctx.Conf.UserAgent)

	// Setup the bootstrap authenticator. Needed to update the site main token.
	ctx.ESIBootstrapAuthenticator = goesi.NewSSOAuthenticator(
		ctx.HTTPClient,
		ctx.Conf.CREST.ESIAccessToken.ClientID,
		ctx.Conf.CREST.ESIAccessToken.SecretKey,
		ctx.Conf.CREST.ESIAccessToken.RedirectURL,
		[]string{"esi-universe.read_structures.v1",
			"esi-search.search_structures.v1",
			"esi-markets.structure_markets.v1"})

	// Get the token from config and build a TokenSource (refreshes the token if needed).
	token := &goesi.CRESTToken{
		AccessToken:  ctx.Conf.CREST.ESIAccessToken.AccessToken,
		TokenType:    ctx.Conf.CREST.ESIAccessToken.TokenType,
		RefreshToken: ctx.Conf.CREST.ESIAccessToken.RefreshToken,
		Expiry:       ctx.Conf.CREST.ESIAccessToken.Expiry,
	}

	ctx.ESIPublicToken, err = ctx.ESIBootstrapAuthenticator.TokenSource(token)
	if err != nil {
		log.Fatalf("Error starting bootstrap ESI client: %v", err)
	}

	// Create a memcached session store.
	ctx.Store, err = gsr.NewRediStoreWithPool(ctx.Cache, []byte(ctx.Conf.Store.Key))
	defer ctx.Store.Close()
	if err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	ctx.Store.SetMaxLength(1024 * 10)
	ctx.Store.Options.Domain = ctx.Conf.Store.Domain

	// Register structs for storage.
	gob.Register(oauth2.Token{})
	gob.Register(goesi.CRESTToken{})
	gob.Register(goesi.VerifyResponse{})

	// Set our logging flags.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if ctx.Conf.Discord.Enabled {
		go discord.GoDiscordBot(&ctx)
	}

	// Run the eve consumers
	if ctx.Conf.EVEConsumer.Enabled {
		eC := eveConsumer.NewEVEConsumer(&ctx)
		eC.RunConsumer()
		defer eC.StopConsumer()
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "dumpdb" {
			err := models.DumpDatabase("./sql/evedata.sql", "evedata")
			if err != nil {
				log.Fatalln(err)
			}
		} else if os.Args[1] == "bootstrap" {
			err := bootstrap(&ctx)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	// Allocate the routes.
	rtr := NewRouter(&ctx)

	log.Printf("EveData Listening port 3000...\n")
	http.ListenAndServe(":3000", context.ClearHandler(rtr))
	log.Printf("EveData Quitting..\n")
}
