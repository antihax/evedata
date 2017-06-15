package zkillboard

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	hammerQueueSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "evedata",
		Subsystem: "hammerqueue",
		Name:      "calls",
		Help:      "Current temperature of the CPU.",
	})
)

func init() {
	prometheus.MustRegister(hammerQueueSize)
}

func (s *ZKillboard) tickMetrics() {
	size, err := s.outQueue.Size()
	if err != nil {
		log.Println(err)
		return
	}
	hammerQueueSize.Set(float64(size))
}
