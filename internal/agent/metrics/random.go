package metrics

import (
	"github.com/tony-spark/metrico/internal/model"
	"math/rand"
	"sync"
	"time"
)

const (
	metricName = "RandomValue"
)

type RandomMetricCollector struct {
	metric  *GaugeMetric
	metrics []model.Metric
	rand    *rand.Rand
	mu      sync.RWMutex
}

func NewRandomMetricCollector() *RandomMetricCollector {
	rmc := &RandomMetricCollector{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	rmc.metric = NewGaugeMetric(metricName, rmc.rand.Float64())
	rmc.metrics = []model.Metric{rmc.metric}
	return rmc
}

func (c *RandomMetricCollector) Metrics() []model.Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.metrics
}

func (c *RandomMetricCollector) Update() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metric.value = c.rand.Float64()
}
