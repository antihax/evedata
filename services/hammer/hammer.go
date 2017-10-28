// Package hammer provides a queued work consumer for CCP ESI API
package hammer

import (
	"log"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	nsq "github.com/nsqio/go-nsq"
)

// Hammer completes work handling CCP ESI and other API.
type Hammer struct {
	stop     chan bool
	hammerWG *sync.WaitGroup
	inQueue  *redisqueue.RedisQueue
	esi      *goesi.APIClient
	redis    *redis.Pool
	nsq      *nsq.Producer
	sem      chan bool

	// authentication
	token *goesi.CRESTTokenSource
	auth  *goesi.SSOAuthenticator
}

// NewHammer Service.
func NewHammer(redis *redis.Pool, nsq *nsq.Producer, clientID string, secret string, refresh string) *Hammer {

	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(redis)

	// Create our ESI API Client
	esi := goesi.NewAPIClient(cache, "EVEData-API-Hammer")

	// Setup an authenticator
	auth := goesi.NewSSOAuthenticator(cache, clientID, secret, "",
		[]string{"esi-universe.read_structures.v1",
			"esi-search.search_structures.v1",
			"esi-markets.structure_markets.v1"})

	tok := &goesi.CRESTToken{
		Expiry:       time.Now(),
		AccessToken:  "",
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}

	// Build our token
	token, err := auth.TokenSource(tok)
	if err != nil {
		log.Fatalln(err)
	}

	// Setup a new hammer
	s := &Hammer{
		stop:     make(chan bool),
		hammerWG: &sync.WaitGroup{},
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		nsq:   nsq,
		auth:  auth,
		esi:   esi,
		redis: redis,
		token: &token,
		sem:   make(chan bool, 50),
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
	s.auth.ChangeTokenURL(path)
	s.auth.ChangeAuthURL(path)
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
