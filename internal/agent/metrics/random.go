package metrics

import (
	"math/rand"
	"time"
)

const (
	metricName = "RandomValue"
)

type RandomMetricCollector struct {
	metric GaugeMetric
	rand   *rand.Rand
}

func NewRandomMetricCollector() *RandomMetricCollector {
	rmc := &RandomMetricCollector{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return rmc
}

func (c *RandomMetricCollector) Metrics() []Metric {
	return []Metric{NewGaugeMetric(metricName, c.rand.Float64())}
}

func (c *RandomMetricCollector) Update() {
	c.metric.value = c.rand.Float64()
}
