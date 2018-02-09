// Package discordbottemp provides a quick hack of a discord bot for fesibility testing.
package discordbottemp

import (
	"log"
	"sync"

	"github.com/antihax/evedata/internal/redisqueue"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
)

// DiscordBot Handles our little bot.
type DiscordBot struct {
	stop              chan bool
	redis             *redis.Pool
	outQueue          *redisqueue.RedisQueue
	db                *sqlx.DB
	discordToken      string
	consumerAddresses []string
	consumers         map[string]*nsq.Consumer
	wg                *sync.WaitGroup
}

// NewDiscordBot Service.
func NewDiscordBot(redis *redis.Pool, db *sqlx.DB, addresses []string, discordToken string) *DiscordBot {
	// Setup a new artifice
	s := &DiscordBot{
		stop:              make(chan bool),
		db:                db,
		redis:             redis,
		discordToken:      discordToken,
		consumerAddresses: addresses,
		consumers:         make(map[string]*nsq.Consumer),
		wg:                &sync.WaitGroup{},
		outQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
	}

	return s
}

// Close the discord service
func (s *DiscordBot) Close() {
	close(s.stop)
	for _, h := range s.consumers {
		h.Stop()
	}
	s.wg.Wait()
}

// Run the discord service
func (s *DiscordBot) Run() {
	err := s.connectToDiscord()
	if err != nil {
		log.Fatal(err)
	}

	err = s.getSystems()
	if err != nil {
		log.Fatal(err)
	}

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
