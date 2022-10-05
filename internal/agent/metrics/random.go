package metrics

import (
	"math/rand"
	"time"
)

const (
	metricName = "RandomValue"
)

type RandomMetricCollector struct {
	metric *GaugeMetric
	rand   *rand.Rand
}

func NewRandomMetricCollector() *RandomMetricCollector {
	rmc := &RandomMetricCollector{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	rmc.metric = NewGaugeMetric(metricName, rmc.rand.Float64())
	return rmc
}

func (c *RandomMetricCollector) Metrics() []Metric {
	return []Metric{c.metric}
}

func (c *RandomMetricCollector) Update() {
	c.metric.value = c.rand.Float64()
}
