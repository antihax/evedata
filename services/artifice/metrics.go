package artifice

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	hammerQueueSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "evedata",
		Subsystem: "hammerqueue",
		Name:      "calls",
		Help:      "Size of the hammer queue.",
	})
)

func init() {
	prometheus.MustRegister(hammerQueueSize)
}

func (s *Artifice) tickMetrics() {
	size, err := s.QueueSize()
	if err != nil {
		log.Println(err)
		return
	}
	hammerQueueSize.Set(float64(size))
}

func (s *Artifice) runMetrics() {
	throttle := time.Tick(time.Second * 5)
	for {
		s.tickMetrics()
		<-throttle
	}
}
