package hammer

import (
	"errors"
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type consumer struct {
	name string
	f    consumerFunc
}

type consumerFunc func(*Hammer, interface{})

// Register a consumer to a queue operation.
func registerConsumer(name string, f consumerFunc) {
	consumers = append(consumers, consumer{name, f})
	consumerMap[name] = f
}

func (s *Hammer) wait(f consumerFunc, p interface{}) {
	// Limit go routines
	s.hammerWG.Add(1)
	defer func() { <-s.sem; s.hammerWG.Done() }()
	f(s, p)
}

func (s *Hammer) runConsumers() error {
	w, err := s.inQueue.GetWork()
	if err != nil {
		return err
	}

	fn := consumerMap[w.Operation]
	if fn == nil {
		log.Printf("unknown operation\n")
		return errors.New("Unknown operation")
	}

	s.sem <- true
	fmt.Printf("running %s %v\n", w.Operation, w.Parameter)
	go s.wait(fn, w.Parameter)

	return nil
}

// For handling Consumers
var (
	consumers       []consumer
	consumerMap     map[string]consumerFunc
	consumerMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "hammer",
		Name:      "ticks",
		Help:      "API call statistics.",
		Buckets:   prometheus.ExponentialBuckets(1, 2, 17),
	}, []string{"consumer"},
	)

	queueSizeMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "evedata",
		Subsystem: "hammer",
		Name:      "size",
		Help:      "Entries in queue for consumers",
	}, []string{"queue"},
	)
)

func init() {
	consumerMap = make(map[string]consumerFunc)

	prometheus.MustRegister(
		consumerMetrics,
		queueSizeMetrics,
	)
}
