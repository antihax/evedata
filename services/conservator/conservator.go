// Package conservator herds service bots for external services.
package conservator

import (
	"log"
	"sync"

	"github.com/antihax/evedata/internal/botservice"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
)

// Conservator Handles our little bot.
type Conservator struct {
	stop              chan bool
	redis             *redis.Pool
	outQueue          *redisqueue.RedisQueue
	db                *sqlx.DB
	consumerAddresses []string
	consumers         map[string]*nsq.Consumer
	services          map[int32]*botservice.BotService
	wg                *sync.WaitGroup
}

// NewConservator Service.
func NewConservator(redis *redis.Pool, db *sqlx.DB, addresses []string) *Conservator {
	// Setup a new artifice
	s := &Conservator{
		stop:  make(chan bool),
		db:    db,
		redis: redis,

		consumerAddresses: addresses,
		consumers:         make(map[string]*nsq.Consumer),
		wg:                &sync.WaitGroup{},
		outQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-test-hammer",
		),
	}

	return s
}

// Close the conservator service
func (s *Conservator) Close() {
	close(s.stop)
	for _, h := range s.consumers {
		h.Stop()
	}
	s.wg.Wait()
}

// Run the conservator service
func (s *Conservator) Run() {
	err := s.getSystems()
	if err != nil {
		log.Fatal(err)
	}

	// Run the war ticker
	go s.updateWars()

	// Run the war ticker
	go s.updateWars()

	err = s.registerHandlers()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-s.stop:
			return
		}
	}
}
