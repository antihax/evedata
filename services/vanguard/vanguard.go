package vanguard

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/discordauth"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/evedata/services/vanguard/models"
	"github.com/antihax/goesi"
	gsr "github.com/antihax/redistore"
	"github.com/coreos/go-oidc/v3/oidc"
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
	Verifier             *oidc.IDTokenVerifier
	DiscordAuthenticator *discordauth.Authenticator
}

var globalVanguard *Vanguard

func NewVanguard(redis *redis.Pool, db *sqlx.DB) *Vanguard {
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
	tauth := goesi.NewSSOAuthenticatorV2(cache, os.Getenv("ESI_CLIENTID_TOKENSTORE"), os.Getenv("ESI_SECRET_TOKENSTORE"), "https://"+os.Getenv("DOMAIN")+"/U/eveTokenAnswer", []string{})

	// Setup an authenticator for our SSO token
	ssoauth := goesi.NewSSOAuthenticatorV2(cache, os.Getenv("ESI_CLIENTID_SSO"), os.Getenv("ESI_SECRET_SSO"), "https://"+os.Getenv("DOMAIN")+"/U/eveSSOAnswer", []string{})

	// Setup an authenticator for Discord
	dauth := discordauth.NewAuthenticator(cache, os.Getenv("DISCORD_CLIENTID"), os.Getenv("DISCORD_SECRET"), "https://"+os.Getenv("DOMAIN")+"/U/discordAnswer", []string{"identify", "guilds.join"})

	// Build our private token
	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: os.Getenv("ESI_REFRESHKEY"),
		TokenType:    "Bearer",
	}
	token := tauth.TokenSource(tok)

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

	key := oidc.NewRemoteKeySet(context.Background(), "https://login.eveonline.com/oauth/jwks")
	if key == nil {
		panic("could not obtain remote key")
	}
	verifier := oidc.NewVerifier("login.eveonline.com", key, &oidc.Config{SkipClientIDCheck: true})

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
		Verifier:             verifier,
		ESI:                  esi,
		Store:                store,
		Db:                   db,
		Cache:                redis,

		token:      &token,
		TokenStore: tokenStore,
	}

	go func(gvg *Vanguard) {
		if err := gvg.RPCConnect(); err != nil {
			log.Fatalln(err)
		}
	}(globalVanguard)

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
	s.Conservator, err = rpc.DialHTTP("tcp", "conservator.evedata.svc.cluster.local:3001")
	if err != nil {
		s.Conservator, err = rpc.DialHTTP("tcp", "conservator.evedata.svc.cluster.local:32003")
		if err != nil {
			return err
		}
	}

	return nil
}

// JWTVerify validates jwt tokens and returns legacy VerifyResponse response
func (s *Vanguard) JWTVerify(token string) (*goesi.VerifyResponse, error) {
	idToken, err := s.Verifier.Verify(context.Background(), token)
	if err != nil {
		return nil, err
	}

	var claims goesi.EVESSOClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	idParts := strings.Split(claims.Subject, ":")
	if len(idParts) != 3 {
		return nil, fmt.Errorf("could not decode character id")
	}

	cid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return nil, err
	}

	v := &goesi.VerifyResponse{
		CharacterName:      claims.Name,
		CharacterID:        int32(cid),
		CharacterOwnerHash: claims.Owner,
	}
	return v, nil
}
