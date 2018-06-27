// Package hammer provides a queued work consumer for CCP ESI API
package hammer

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/oauth2"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	nsq "github.com/nsqio/go-nsq"
)

const NUM_WORKERS = 100

// Hammer completes work handling CCP ESI and other API.
type Hammer struct {
	stop       chan bool
	wg         *sync.WaitGroup
	inQueue    *redisqueue.RedisQueue
	esi        *goesi.APIClient
	db         *sqlx.DB
	redis      *redis.Pool
	nsq        *nsq.Producer
	sem        chan bool
	tokenStore *tokenstore.TokenStore

	// Count of active workers
	activeWorkers uint64

	// authentication
	token     *oauth2.TokenSource
	tokenAuth *goesi.SSOAuthenticator
}

// NewHammer Service.
func NewHammer(redis *redis.Pool, ledis *redis.Pool, db *sqlx.DB, nsq *nsq.Producer, refresh, tokenClientID, tokenSecret string) *Hammer {
	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(ledis)

	// Create our ESI API Client
	esi := goesi.NewAPIClient(cache, "EVEData-API-Hammer")

	// Setup an authenticator for our user tokens
	tauth := goesi.NewSSOAuthenticator(cache, tokenClientID, tokenSecret, "", []string{})

	tok := &oauth2.Token{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}

	// Build our private token
	token, err := tauth.TokenSource(tok)
	if err != nil {
		log.Fatalln(err)
	}

	tokenStore := tokenstore.NewTokenStore(redis, db, tauth)

	// Setup a new hammer
	s := &Hammer{
		stop: make(chan bool),
		wg:   &sync.WaitGroup{},
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		nsq:        nsq,
		tokenAuth:  tauth,
		esi:        esi,
		db:         db,
		redis:      redis,
		token:      &token,
		tokenStore: tokenStore,
		sem:        make(chan bool, NUM_WORKERS),
	}

	return s
}

// Close the hammer service
func (s *Hammer) Close() {
	close(s.stop)
	s.wg.Wait()
}

// ChangeBasePath for ESI (sisi/mock/tranquility)
func (s *Hammer) ChangeBasePath(path string) {
	s.esi.ChangeBasePath(path)
}

// ChangeTokenPath for ESI (sisi/mock/tranquility)
func (s *Hammer) ChangeTokenPath(path string) {
	s.tokenAuth.ChangeTokenURL(path)
	s.tokenAuth.ChangeAuthURL(path)
}

// QueueWork directly
func (s *Hammer) QueueWork(work []redisqueue.Work, priority int) error {
	return s.inQueue.QueueWork(work, priority)
}

// Run the service
func (s *Hammer) Run() {
	go s.tickWorkersToPrometheus()

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

func (s *Hammer) tickWorkersToPrometheus() {
	for {
		<-time.After(5 * time.Second)
		count := atomic.LoadUint64(&s.activeWorkers)
		activeWorkers.Set(float64(count))
	}
}

var (
	activeWorkers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "evedata",
			Subsystem: "hammer",
			Name:      "activeWorkers",
			Help:      "Currently running workers.",
		},
	)
)

func init() {
	prometheus.MustRegister(
		activeWorkers,
	)
}
