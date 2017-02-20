package eveConsumer

import (
	"log"
	"strings"
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
	name      string
	f         consumerFunc
	queueName string
}
type consumerFunc func(*EVEConsumer, *redis.Conn) (bool, error)

func addConsumer(name string, f consumerFunc, queueName string) {
	consumers = append(consumers, consumer{name, f, queueName})
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

	queueSizeMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "evedata",
		Subsystem: "queue",
		Name:      "size",
		Help:      "Entries in queue for consumers",
	}, []string{"queue"},
	)
)

func init() {
	prometheus.MustRegister(
		consumerMetrics,
		triggerMetrics,
		queueSizeMetrics,
	)
}

// EVEConsumer provides the microservice which conducts backend
// polling of EVE Crest and XML resources as needed.
type EVEConsumer struct {
	ctx                 *appContext.AppContext
	consumerStopChannel chan bool
	triggersStopChannel chan bool
	metricsStopChannel  chan bool
	errorRate           int32
}

// NewEVEConsumer creates a new EveConsumer
func NewEVEConsumer(ctx *appContext.AppContext) *EVEConsumer {
	e := &EVEConsumer{ctx, make(chan bool), make(chan bool), make(chan bool), 0}
	return e
}

func (c *EVEConsumer) goMetrics() {
	rate := time.Second * 5
	throttle := time.Tick(rate)

	// Run Phase
	for {
		<-throttle

		select {
		case <-c.metricsStopChannel:
			return
		default:
			r := c.ctx.Cache.Get()
			for _, consumer := range consumers {
				if consumer.queueName != "" {
					v, err := redis.Int(r.Do("SCARD", consumer.queueName))
					if err != nil {
						log.Printf("%s: %v\n", consumer.queueName, err)
						continue
					}

					queueName := strings.Replace(consumer.queueName, "EVEDATA_", "", 1)
					queueName = strings.Replace(queueName, "Queue", "", 1)

					queueSizeMetrics.With(
						prometheus.Labels{"queue": queueName},
					).Set(float64(v))
				}
			}
			r.Close()
		}
	}
}

func (c *EVEConsumer) goConsumer() {
	// Run Phase
	for {

		workDone := false
		select {
		case <-c.consumerStopChannel:
			return
		default:
			r := c.ctx.Cache.Get()
			// loop through all the consumers
			for _, consumer := range consumers {
				start := monotime.Now()
				// Call the function
				if workDone, err := consumer.f(c, &r); err == nil {
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
			r.Close()
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
	c.initKillConsumer()
}

// RunConsumer starts the consumer and returns.
func (c *EVEConsumer) RunConsumer() {
	// Load deferrable data.
	go c.initConsumer()
	go c.goMetrics()

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
	c.metricsStopChannel <- true
	c.triggersStopChannel <- true
	log.Printf("EVEConsumer: Stopped\n")
}
