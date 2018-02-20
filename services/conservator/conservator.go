// Package conservator herds service bots for external services.
package conservator

import (
	"log"
	"sync"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/bwmarrin/discordgo"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
)

var NOTIFICATION_TYPES = []string{"kill", "war", "locator", "structure"}

// Conservator Handles our little bot.
type Conservator struct {
	stop              chan bool
	redis             *redis.Pool
	outQueue          *redisqueue.RedisQueue
	db                *sqlx.DB
	consumerAddresses []string
	consumers         map[string]*nsq.Consumer

	wg           *sync.WaitGroup
	discord      *discordgo.Session
	discordToken string

	solarSystems map[int32]float32

	// Base Data
	services sync.Map
	channels sync.Map
	shares   sync.Map
	warsMap  map[int32]*sync.Map

	// Notification data
	notifications    map[string]map[int32][]Share
	notificationLock map[string]*sync.RWMutex
}

// NewConservator Service.
func NewConservator(redis *redis.Pool, db *sqlx.DB, addresses []string, discordToken string) *Conservator {
	// Setup a new artifice
	s := &Conservator{
		stop:              make(chan bool),
		db:                db,
		redis:             redis,
		discordToken:      discordToken,
		consumerAddresses: addresses,
		consumers:         make(map[string]*nsq.Consumer),
		warsMap:           make(map[int32]*sync.Map),
		wg:                &sync.WaitGroup{},
		outQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
	}

	s.notifications = make(map[string]map[int32][]Share)
	s.notificationLock = make(map[string]*sync.RWMutex)
	for _, t := range NOTIFICATION_TYPES {
		s.notifications[t] = make(map[int32][]Share)
		s.notificationLock[t] = &sync.RWMutex{}
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
	var err error

	s.discord, err = discordgo.New("Bot " + s.discordToken)

	// Load solarSystems
	s.solarSystems, err = s.getSolarSystems()
	if err != nil {
		log.Fatal(err)
	}

	// Run the API
	err = s.runRPC()
	if err != nil {
		log.Fatal(err)
	}

	if err = s.loadServices(); err != nil {
		log.Fatal(err)
	}

	if err = s.loadChannels(); err != nil {
		log.Fatal(err)
	}

	if err = s.loadShares(); err != nil {
		log.Fatal(err)
	}

	if err = s.registerHandlers(); err != nil {
		log.Fatal(err)
	}

	// update data
	go s.updateData()

	// Loop until we stop
	for {
		select {
		case <-s.stop:
			return
		}
	}
}

func inSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
