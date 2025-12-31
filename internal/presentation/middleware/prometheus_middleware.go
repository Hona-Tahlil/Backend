package middleware

import (
	domainmetrics "hona/backend/internal/domain/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PrometheusMiddleware struct {
	metricsClient domainmetrics.PrometheusMetrics
}

func NewPrometheusMiddleware(metricsClient domainmetrics.PrometheusMetrics) *PrometheusMiddleware {
	return &PrometheusMiddleware{
		metricsClient: metricsClient,
	}
}

func (pm *PrometheusMiddleware) PrometheusMiddleware(ctx *gin.Context) {
	start := time.Now()

	ctx.Next()

	duration := time.Since(start).Seconds()
	status := ctx.Writer.Status()

	route := ctx.Request.URL.Path
	method := ctx.Request.Method

	pm.metricsClient.IncHTTPRequest(method, route, strconv.Itoa(status))
	pm.metricsClient.ObserveHTTPRequestDuration(method, route, duration)
}
