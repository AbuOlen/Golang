package telemetry

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"module6/internal/config"
	"strconv"
	"time"
)

var reqLatencyHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "request_latency",
	Help: "Latency of HTTP requests",
	Buckets: []float64{
		float64(time.Second) * .005,
		float64(time.Second) * .01,
		float64(time.Second) * .025,
		float64(time.Second) * .05,
		float64(time.Second) * .1,
		float64(time.Second) * .25,
		float64(time.Second) * .5,
		float64(time.Second) * 1,
		float64(time.Second) * 2.5,
		float64(time.Second) * 5,
		float64(time.Second) * 10,
	},
}, []string{
	"instance",
	"path",
	"status",
})

type Middleware struct {
	cfg *config.Config
}

func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{
		cfg: cfg,
	}
}

func (m *Middleware) Handle(ctx *fiber.Ctx) error {
	startedAt := time.Now()

	err := ctx.Next()

	delay := time.Since(startedAt)
	labels := prometheus.Labels{
		"instance": m.cfg.InstanceName,
		"path":     ctx.Route().Path,
		"status":   strconv.Itoa(ctx.Response().StatusCode()),
	}

	reqLatencyHistogram.With(labels).Observe(float64(delay))

	return err
}
