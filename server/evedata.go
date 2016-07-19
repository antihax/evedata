package evedata

import (
	"database/sql"
	"encoding/gob"
	"evedata/config"
	"evedata/eveapi"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/bradfitz/gomemcache/memcache"
	gsm "github.com/bradleypeabody/gorilla-sessions-memcache"
	"github.com/gorilla/context"
	"github.com/gregjones/httpcache"
	httpmemcache "github.com/gregjones/httpcache/memcache" // dumb package name...

	"github.com/jmoiron/sqlx"
)

// appContext provides access to handles throughout the app.
type AppContext struct {
	Conf  *config.Config
	Db    *sqlx.DB
	Store *gsm.MemcacheStore

	SSOAuthenticator *eveapi.SSOAuthenticator

	HTTPClient *http.Client

	Bridge struct {
		HistoryUpdate *sql.Stmt
		OrderMark     *sql.Stmt
		OrderUpdate   *sql.Stmt
		KillInsert    *sql.Stmt
	}
}

func GoServer() {
	var err error

	// Make a new app context.8
	ctx := &AppContext{}

	// Read configuation.
	ctx.Conf, err = config.ReadConfig()

	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Build Connection Pool
	ctx.Db, err = sqlx.Connect(ctx.Conf.Database.Driver, ctx.Conf.Database.Spec)
	if err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	// Check we can connect
	err = ctx.Db.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	ctx.SSOAuthenticator = eveapi.NewSSOAuthenticator(ctx.Conf.CREST.ClientID,
		ctx.Conf.CREST.SecretKey,
		ctx.Conf.CREST.RedirectURL)

	// Allocate the routes
	rtr := NewRouter(ctx)

	// Connect to the memcache server
	cache := memcache.New(ctx.Conf.MemcachedAddress)

	// Create a memcached http client for the CCP APIs.
	crestCache := httpmemcache.NewWithClient(cache)
	transport := httpcache.NewTransport(crestCache)
	ctx.HTTPClient = &http.Client{Transport: transport}

	// Create a memcached session store.
	ctx.Store = gsm.NewMemcacheStore(cache, "EVEDATA_SESSIONS_", []byte(ctx.Conf.Store.Key))

	gob.Register(eveapi.VerifyResponse{})
	gob.Register(oauth2.Token{})
	ctx.Store.StoreMethod = gsm.StoreMethodSecureCookie

	// Set our logging flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if ctx.Conf.EMDRCrestBridge.Enabled {
		log.Println("Starting EMDR <- Crest Bridge")
		go goEMDRCrestBridge(ctx)
	}

	log.Printf("EveData Listening port 3000...\n")
	http.ListenAndServe(":3000", context.ClearHandler(rtr))
}
