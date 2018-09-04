package tailor

import "github.com/prometheus/client_golang/prometheus"

var (
	metricStageTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "evedata",
		Subsystem: "tailor",
		Name:      "stageTime",
		Help:      "Stage Durations",
		Buckets:   prometheus.ExponentialBuckets(1, 1.2, 50),
	},
		[]string{"stage"},
	)

	metricStatus = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "evedata",
		Subsystem: "tailor",
		Name:      "status",
		Help:      "Count of tailor uploads.",
	},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(
		metricStageTime,
		metricStatus,
	)
}
