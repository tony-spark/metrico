package transports

import "github.com/tony-spark/metrico/internal/agent/metrics"

type Transport interface {
	SendMetric(metric metrics.Metric) error
	SendMetrics(mx []metrics.Metric) error
}
