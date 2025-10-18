package integration

import "github.com/prometheus/client_golang/prometheus"

// MetricsRegistry stores Prometheus collectors for integration features.
type MetricsRegistry struct {
	collectors []prometheus.Collector
}

// NewMetricsRegistry creates an empty metrics registry ready for collectors.
func NewMetricsRegistry() *MetricsRegistry {
	return &MetricsRegistry{}
}

// Collectors returns the registered Prometheus collectors.
func (r *MetricsRegistry) Collectors() []prometheus.Collector {
	return r.collectors
}
