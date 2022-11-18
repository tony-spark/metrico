package metrics

import (
	"fmt"
	"github.com/tony-spark/metrico/internal/model"
)

type MetricCollector interface {
	Metrics() []model.Metric
	Update()
}

type GaugeMetric struct {
	name  string
	value float64
}

type CounterMetric struct {
	name  string
	value int64
}

func (g GaugeMetric) String() string {
	return fmt.Sprint(g.value)
}

func (g GaugeMetric) ID() string {
	return g.name
}

func (g GaugeMetric) Type() string {
	return model.GAUGE
}

func (g GaugeMetric) Val() interface{} {
	return g.value
}

func (c CounterMetric) String() string {
	return fmt.Sprint(c.value)
}

func (c CounterMetric) ID() string {
	return c.name
}

func (c CounterMetric) Type() string {
	return model.COUNTER
}

func (c CounterMetric) Val() interface{} {
	return c.value
}

func NewGaugeMetric(name string, value float64) *GaugeMetric {
	return &GaugeMetric{name, value}
}

func NewCounterMetric(name string, value int64) *CounterMetric {
	return &CounterMetric{name, value}
}
