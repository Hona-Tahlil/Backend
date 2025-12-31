package domainmetrics

type PrometheusMetrics interface {
	IncHTTPRequest(method, route, status string)
	ObserveHTTPRequestDuration(method, route string, duration float64)
}
