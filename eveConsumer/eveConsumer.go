package eveConsumer

import (
	"log"
	"strings"
	"time"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/internal/redisqueue"
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
		Buckets:   prometheus.ExponentialBuckets(10, 1.45, 20),
	}, []string{"consumer"},
	)

	triggers       []trigger
	triggerMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "trigger",
		Name:      "ticks",
		Help:      "API call statistics.",
		Buckets:   prometheus.ExponentialBuckets(10, 1.45, 20),
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

var hammerQueue *redisqueue.RedisQueue

// NewEVEConsumer creates a new EveConsumer
func NewEVEConsumer(ctx *appContext.AppContext) *EVEConsumer {
	hammerQueue = redisqueue.NewRedisQueue(ctx.Cache, "evedata-hammer")
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
		var (
			err      error
			workDone bool
		)

		select {
		case <-c.consumerStopChannel:
			return
		default:
			r := c.ctx.Cache.Get()
			// loop through all the consumers
			for _, consumer := range consumers {
				start := time.Now()
				// Call the function
				if workDone, err = consumer.f(c, &r); err == nil {
					if workDone {
						duration := float64(time.Since(start).Nanoseconds()) / 1000000.0
						consumerMetrics.With(
							prometheus.Labels{"consumer": consumer.name},
						).Observe(duration)
					}
				} else if err != nil {
					log.Printf("%s: %v\n", consumer.name, err)
				}
			}
			r.Close()
		}

		// Sleep a brief bit if we didnt do anything
		if workDone == false {
			time.Sleep(time.Second)
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
				start := time.Now()
				if workDone, err := trigger.f(c); err == nil {
					if workDone {
						duration := float64(time.Since(start).Nanoseconds()) / 1000000.0
						triggerMetrics.With(
							prometheus.Labels{"trigger": trigger.name},
						).Observe(duration)
					}
				} else if err != nil {
					log.Printf("%s: %v\n", trigger.name, err)
				}
			}
		}
		<-throttle
	}
}

// RunConsumer starts the consumer and returns.
func (c *EVEConsumer) RunConsumer() {
	// Load deferrable data.
	go c.goMetrics()

	for i := 0; i < c.ctx.Conf.EVEConsumer.Consumers; i++ {
		go c.goConsumer()                 // Run consumers in a loop
		time.Sleep(time.Millisecond * 37) // Stagger starting the routines
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
