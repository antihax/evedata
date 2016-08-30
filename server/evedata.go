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

	"golang.org/x/oauth2"

	"github.com/bradfitz/gomemcache/memcache"
	gsm "github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/gorilla/context"
	"github.com/gregjones/httpcache"
	httpmemcache "github.com/gregjones/httpcache/memcache" // ...
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

	// Connect to the memcache server
	cache := memcache.New(ctx.Conf.MemcachedAddress)

	// Create a memcached http client for the CCP APIs.
	transport := httpcache.NewTransport(httpmemcache.NewWithClient(cache))
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	ctx.HTTPClient = &http.Client{Transport: transport}

	// Create a memcached session store.
	ctx.Store = gsm.NewMemcacheStore(cache, "EVEDATA_SESSIONS_", []byte(ctx.Conf.Store.Key))
	ctx.Store.StoreMethod = gsm.StoreMethodSecureCookie
	ctx.Store.Options.Domain = ctx.Conf.Domain

	// Register structs for storage
	gob.Register(oauth2.Token{})
	gob.Register(eveapi.CRESTToken{})

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
