package appContext

import (
	"evedata/config"
	"evedata/esi"
	"evedata/eveapi"
	"net/http"

	"golang.org/x/oauth2"

	redistore "gopkg.in/boj/redistore.v1"

	"github.com/garyburd/redigo/redis"
	"github.com/gregjones/httpcache"
	"github.com/jmoiron/sqlx"
)

// AppContext provides access to handles throughout the app.
type AppContext struct {
	Conf           *config.Config          // App Configuration
	Db             *sqlx.DB                // EVE Database
	Store          *redistore.RediStore    // Redis session store.
	EVE            *eveapi.AnonymousClient // EVE API Client
	HTTPClient     *http.Client            // Redis Cached HTTP client
	TransportCache *httpcache.Transport
	Cache          redis.Pool // Redis connection Pool for HTTP Cache and session store.
	ESI            *esi.APIClient
	ESIPublicToken oauth2.TokenSource

	// Since we need to combine data from multiple characters, we use
	// one authenticator for the site to act as the main authentication.
	// second will allow for many alt characters under the main.
	SSOAuthenticator          *eveapi.SSOAuthenticator // CREST authenticator for site authentication
	TokenAuthenticator        *eveapi.SSOAuthenticator // CREST authenticator for site functionality
	ESIBootstrapAuthenticator *eveapi.SSOAuthenticator // CREST authenticator for site functionality
}
