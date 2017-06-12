// Package hammer provides a queued work consumer for CCP ESI API
package hammer

import (
	"log"
	"sync"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/goesi"
	"github.com/garyburd/redigo/redis"
	nsq "github.com/nsqio/go-nsq"
)

// Hammer provides service control.
type Hammer struct {
	stop     chan bool
	hammerWG *sync.WaitGroup
	inQueue  *redisqueue.RedisQueue
	esi      *goesi.APIClient
	redis    *redis.Pool
	nsq      *nsq.Producer
}

// NewHammer Service.
func NewHammer(redis *redis.Pool, nsq *nsq.Producer) *Hammer {
	// Get a caching http client
	cache := apicache.CreateHTTPClientCache(redis)

	// Create our ESI API Client
	esi := goesi.NewAPIClient(cache, "EVEData-API-Hammer")

	// Setup a new hammer
	s := &Hammer{
		stop:     make(chan bool),
		hammerWG: &sync.WaitGroup{},
		inQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		nsq:   nsq,
		esi:   esi,
		redis: redis,
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
