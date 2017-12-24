// Package hammer provides a queued work consumer for CCP ESI API
package hammer

import (
	"context"
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	nsq "github.com/nsqio/go-nsq"
)

// Hammer completes work handling CCP ESI and other API.
type Hammer struct {
	stop       chan bool
	hammerWG   *sync.WaitGroup
	inQueue    *redisqueue.RedisQueue
	esi        *goesi.APIClient
	db         *sqlx.DB
	redis      *redis.Pool
	nsq        *nsq.Producer
	sem        chan bool
	tokenStore *tokenstore.TokenStore

	// authentication
	token       *oauth2.TokenSource
	privateAuth *goesi.SSOAuthenticator
	tokenAuth   *goesi.SSOAuthenticator
}

// NewHammer Service.
func NewHammer(redis *redis.Pool, db *sqlx.DB, nsq *nsq.Producer, clientID, secret, refresh, tokenClientID, tokenSecret string) *Hammer {
	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(redis)

	// Create our ESI API Client
	esi := goesi.NewAPIClient(cache, "EVEData-API-Hammer")

	// Setup an authenticator for our user tokens
	tauth := goesi.NewSSOAuthenticator(cache, tokenClientID, tokenSecret, "", []string{})

	// Setup an authenticator for our private token
	pauth := goesi.NewSSOAuthenticator(cache, clientID, secret, "",
		[]string{"esi-universe.read_structures.v1",
			"esi-search.search_structures.v1",
			"esi-markets.structure_markets.v1"})

	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}

	// Build our private token
	token, err := pauth.TokenSource(tok)
	if err != nil {
		log.Fatalln(err)
	}

	tokenStore := tokenstore.NewTokenStore(redis, db, tauth)

	// Setup a new hammer
	s := &Hammer{
		stop:     make(chan bool),
		hammerWG: &sync.WaitGroup{},
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		nsq:         nsq,
		privateAuth: pauth,
		tokenAuth:   tauth,
		esi:         esi,
		db:          db,
		redis:       redis,
		token:       &token,
		tokenStore:  tokenStore,
		sem:         make(chan bool, 100),
	}

	return s
}

// Close the hammer service
func (s *Hammer) Close() {
	close(s.stop)
	s.hammerWG.Wait()
}

// ChangeBasePath for ESI (sisi/mock/tranquility)
func (s *Hammer) ChangeBasePath(path string) {
	s.esi.ChangeBasePath(path)
}

// ChangeTokenPath for ESI (sisi/mock/tranquility)
func (s *Hammer) ChangeTokenPath(path string) {
	s.privateAuth.ChangeTokenURL(path)
	s.privateAuth.ChangeAuthURL(path)
	s.tokenAuth.ChangeTokenURL(path)
	s.tokenAuth.ChangeAuthURL(path)
}

// QueueWork directly
func (s *Hammer) QueueWork(work []redisqueue.Work) error {
	return s.inQueue.QueueWork(work)
}

// Run the hammer service
func (s *Hammer) Run() {
	for {
		select {
		case <-s.stop:
			return
		default:
			err := s.runConsumers()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// GetTokenSourceContext sets a token source to a context value for authentication
func (s *Hammer) GetTokenSourceContext(c context.Context, characterID, tokenCharacterID int32) (context.Context, error) {
	tokenSource, err := s.tokenStore.GetTokenSource(characterID, tokenCharacterID)
	if err != nil {
		return c, err
	}
	auth := context.WithValue(c, goesi.ContextOAuth2, tokenSource)
	return auth, nil
}

// QueueResult queues a result to NSQ topic
func (s *Hammer) QueueResult(v interface{}, topic string) error {
	b, err := bson.Marshal(v)
	if err != nil {
		return err
	}

	return s.nsq.Publish(topic, b)
}

// SetToken Sets a token to the store
func (s *Hammer) SetToken(cid, tcid int32, token *oauth2.Token) error {
	return s.tokenStore.SetToken(cid, tcid, token)
}
