package hammer

import (
	"errors"
	"log"
	"sync/atomic"
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
	s.wg.Add(1)
	atomic.AddUint64(&s.activeWorkers, 1)
	defer func() { <-s.sem; s.wg.Done(); atomic.AddUint64(&s.activeWorkers, ^uint64(0)) }()
	f(s, p)
}

func (s *Hammer) runConsumers() error {
	w, err := s.inQueue.GetWork()
	if err != nil {
		return err
	}

	start := time.Now()
	fn := consumerMap[w.Operation]
	if fn == nil {
		log.Printf("unknown operation %s %+v\n", w.Operation, w.Parameter)
		return errors.New("Unknown operation")
	}

	s.sem <- true
	go s.wait(fn, w.Parameter)

	duration := float64(time.Since(start).Nanoseconds()) / 1000000.0
	consumerMetrics.With(
		prometheus.Labels{"operation": w.Operation},
	).Observe(duration)

	return nil
}

var (
	consumerMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "hammer",
		Name:      "ticks",
		Help:      "API call statistics.",
		Buckets:   prometheus.ExponentialBuckets(10, 1.45, 20),
	}, []string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(
		consumerMetrics,
	)
}
