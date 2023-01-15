// Package transports contains various implementations of sending metrics to server
package transports

import (
	"context"

	"github.com/tony-spark/metrico/internal/model"
)

type Transport interface {
	SendMetric(metric model.Metric) error
	SendMetrics(mx []model.Metric) error
	SendMetricsWithContext(ctx context.Context, mx []model.Metric) error
}
