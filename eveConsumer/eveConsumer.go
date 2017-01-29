package eveConsumer

import (
	"log"
	"time"

	"github.com/ScaleFT/monotime"
	"github.com/antihax/evedata/appContext"
	"github.com/garyburd/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"
)

// For handling triggers
type triggerFunc func(*EVEConsumer) (bool, error)
type trigger struct {
	name string
	f    triggerFunc
}

func addTrigger(name string, f triggerFunc) {
	triggers = append(triggers, trigger{name, f})
}

type consumer struct {
	name string
	f    consumerFunc
}
type consumerFunc func(*EVEConsumer, redis.Conn) (bool, error)

func addConsumer(name string, f consumerFunc) {
	consumers = append(consumers, consumer{name, f})

}

// For handling Consumers
var (
	consumers       []consumer
	consumerMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "consumer",
		Name:      "ticks",
		Help:      "API call statistics.",
		Buckets:   prometheus.ExponentialBuckets(1, 2, 17),
	}, []string{"consumer"},
	)

	triggers       []trigger
	triggerMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "trigger",
		Name:      "ticks",
		Help:      "API call statistics.",
		Buckets:   prometheus.ExponentialBuckets(1, 2, 17),
	}, []string{"trigger"},
	)
)

func init() {
	prometheus.MustRegister(
		consumerMetrics,
		triggerMetrics,
	)
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
				start := monotime.Now()
				if workDone, err := consumer.f(c, r); err == nil {
					if workDone {
						duration := monotime.Duration(start, monotime.Now())
						consumerMetrics.With(
							prometheus.Labels{"consumer": consumer.name},
						).Observe(float64(duration / time.Millisecond))
					}
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
	// Run this every 5 minutes.
	// The triggers should have their own internal checks for cache timers
	rate := time.Second * 60 * 1
	throttle := time.Tick(rate)
	for {
		select {
		case <-c.triggersStopChannel:
			log.Printf("EVEConsumer: Triggers shutting down\n")
			return
		default:
			// loop through all the consumers
			for _, trigger := range triggers {
				start := monotime.Now()
				if workDone, err := trigger.f(c); err == nil {
					if workDone {
						duration := monotime.Duration(start, monotime.Now())
						triggerMetrics.With(
							prometheus.Labels{"trigger": trigger.name},
						).Observe(float64(duration / time.Millisecond))
					}
				} else if err != nil {
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
