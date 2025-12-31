package metrics

import (
	"hona/backend/bootstrap"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusMetrics struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{

				Name: bootstrap.Run().Constants.Metrics.HTTPRequestsTotal.Name,
				Help: bootstrap.Run().Constants.Metrics.HTTPRequestsTotal.Help,
			},
			[]string{"method", "route", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    bootstrap.Run().Constants.Metrics.HTTPRequestDuration.Name,
				Help:    bootstrap.Run().Constants.Metrics.HTTPRequestDuration.Help,
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route"},
		),
	}
}

func (pm *PrometheusMetrics) IncHTTPRequest(method, route, status string) {
	pm.httpRequestsTotal.WithLabelValues(method, route, status).Inc()
}

func (pm *PrometheusMetrics) ObserveHTTPRequestDuration(method, route string, duration float64) {
	pm.httpRequestDuration.WithLabelValues(method, route).Observe(duration)
}
