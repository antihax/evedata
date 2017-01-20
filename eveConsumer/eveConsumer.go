package eveConsumer

import (
	"log"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/garyburd/redigo/redis"
)

// For handling triggers
type triggerFunc func(*EVEConsumer) error
type trigger struct {
	name string
	f    triggerFunc
}

var triggers []trigger

func addTrigger(name string, f triggerFunc) {
	triggers = append(triggers, trigger{name, f})
}

// For handling Consumers
var consumers []consumer

type consumer struct {
	name string
	f    consumerFunc
}
type consumerFunc func(*EVEConsumer, redis.Conn) error

func addConsumer(name string, f consumerFunc) {
	consumers = append(consumers, consumer{name, f})
}

// EVEConsumer provides the microservice which conducts backend
// polling of EVE Crest and XML resources as needed.
type EVEConsumer struct {
	ctx                 *appContext.AppContext
	consumerStopChannel chan bool
	triggersStopChannel chan bool
}

// NewEVEConsumer creates a new EveConsumer
func NewEVEConsumer(ctx *appContext.AppContext) *EVEConsumer {
	e := &EVEConsumer{ctx, make(chan bool), make(chan bool)}
	return e
}

func (c *EVEConsumer) goConsumer() {
	r := c.ctx.Cache.Get()
	defer r.Close()

	// Run Phase
	for {
		workDone := false
		select {
		case <-c.consumerStopChannel:
			return
		default:
			// loop through all the consumers
			for _, consumer := range consumers {
				if err := consumer.f(c, r); err == nil {
					workDone = true
				} else if err != nil {
					log.Printf("%s: %v\n", consumer.name, err)
				}
			}
		}

		// Sleep a brief bit if we didnt do anything
		if workDone == false {
			time.Sleep(time.Second * 5)
		}
	}
}

func (c *EVEConsumer) goTriggers() {
	log.Printf("EVEConsumer: Running Triggers\n")
	// Run this every 5 minutes.
	// The triggers should have their own internal checks for cache timers
	rate := time.Second * 60 * 5
	throttle := time.Tick(rate)
	for {
		select {
		case <-c.triggersStopChannel:
			log.Printf("EVEConsumer: Triggers shutting down\n")
			return
		default:
			// loop through all the consumers
			for _, trigger := range triggers {
				if err := trigger.f(c); err != nil {
					log.Printf("%s: %v\n", trigger.name, err)
				}
			}
		}
		<-throttle
	}
}

// Load deferrable data.
func (c *EVEConsumer) initConsumer() {
	r := c.ctx.Cache.Get()
	defer r.Close()
	// Load Phase
	c.initKillConsumer(r)
}

// RunConsumer starts the consumer and returns.
func (c *EVEConsumer) RunConsumer() {
	// Load deferrable data.
	go c.initConsumer()

	for i := 0; i < c.ctx.Conf.EVEConsumer.Consumers; i++ {
		go c.goConsumer() // Run consumers in a loop
	}

	go c.goTriggers() // Time triggered queries
	if c.ctx.Conf.EVEConsumer.ZKillEnabled == true {
		go c.goZKillConsumer()
		go c.goZKillTemporaryConsumer()
	}

	log.Printf("EVEConsumer: Started\n")
}

// StopConsumer shuts down any running go routines and returns.
func (c *EVEConsumer) StopConsumer() {
	log.Printf("EVEConsumer: Stopping Consumer\n")
	for i := 0; i > c.ctx.Conf.EVEConsumer.Consumers; i++ {
		c.consumerStopChannel <- true
	}
	c.triggersStopChannel <- true
	log.Printf("EVEConsumer: Stopped\n")
}
