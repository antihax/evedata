package hammer

import (
	"errors"
	"log"
	"time"

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

	start := time.Now().Nanosecond()

	fn := consumerMap[w.Operation]
	if fn == nil {
		log.Printf("unknown operation %s %+v\n", w.Operation, w.Parameter)
		return errors.New("Unknown operation")
	}

	s.sem <- true
	go s.wait(fn, w.Parameter)
	log.Printf("operation complete %s %+v\n", w.Operation, w.Parameter)
	duration := (time.Now().Nanosecond() - start) / 1000000
	consumerMetrics.With(
		prometheus.Labels{"operation": w.Operation},
	).Observe(float64(duration))

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
	}, []string{"operation"},
	)
)

func init() {
	consumerMap = make(map[string]consumerFunc)

	prometheus.MustRegister(
		consumerMetrics,
	)
}
