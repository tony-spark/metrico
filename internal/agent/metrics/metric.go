package metrics

import (
	"fmt"
	"github.com/tony-spark/metrico/internal"
)

type Metric interface {
	fmt.Stringer
	Name() string
	Type() string
}

type MetricCollector interface {
	Metrics() []Metric
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

func (g GaugeMetric) Name() string {
	return g.name
}

func (g GaugeMetric) Type() string {
	return internal.GAUGE
}

func (c CounterMetric) String() string {
	return fmt.Sprint(c.value)
}

func (c CounterMetric) Name() string {
	return c.name
}

func (c CounterMetric) Type() string {
	return internal.COUNTER
}

func NewGaugeMetric(name string, value float64) *GaugeMetric {
	return &GaugeMetric{name, value}
}

func NewCounterMetric(name string, value int64) *CounterMetric {
	return &CounterMetric{name, value}
}
