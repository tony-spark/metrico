package transports

import (
	"github.com/tony-spark/metrico/internal/model"
)

type Transport interface {
	SendMetric(metric model.Metric) error
	SendMetrics(mx []model.Metric) error
}
