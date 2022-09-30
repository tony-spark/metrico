package metrics

import (
	"fmt"
	"github.com/tony-spark/metrico/internal"
	"math/rand"
	"time"
)

type RandomMetric struct {
	v float64
}

type RandomMetricCollector struct {
	metric RandomMetric
	rand   *rand.Rand
}

func NewRandomMetricCollector() *RandomMetricCollector {
	rmc := &RandomMetricCollector{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return rmc
}

func (r RandomMetric) String() string {
	return fmt.Sprint(r.v)
}

func (r RandomMetric) Name() string {
	return "RandomValue"
}

func (r RandomMetric) Type() string {
	return internal.GAUGE
}

func (c *RandomMetricCollector) Metrics() []Metric {
	return []Metric{c.metric}
}

func (c *RandomMetricCollector) Update() {
	c.metric.v = c.rand.Float64()
}
