package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	authLatencyHistogram *prometheus.HistogramVec
)

func init() {
	authLatencyHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "token_latency",
		Help: "Latency of parsing JWT token",
		Buckets: []float64{
			float64(time.Millisecond) * 1,
			float64(time.Millisecond) * 2.5,
			float64(time.Millisecond) * 5,
			float64(time.Millisecond) * 10,
		},
	}, []string{
		"instance",
	})
}

func MeasureAuthTime(instance string, delay time.Duration) {
	authLatencyHistogram.WithLabelValues(instance).Observe(float64(delay.Milliseconds()))
}
