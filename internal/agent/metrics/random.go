package metrics

import (
	"math/rand"
	"sync"
	"time"

	"github.com/tony-spark/metrico/internal/model"
)

const (
	metricName = "RandomValue"
)

type RandomMetricCollector struct {
	metric *GaugeMetric
	rand   *rand.Rand
	mu     sync.RWMutex
}

func NewRandomMetricCollector() *RandomMetricCollector {
	rmc := &RandomMetricCollector{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	rmc.metric = NewGaugeMetric(metricName, rmc.rand.Float64())
	return rmc
}

func (c *RandomMetricCollector) Metrics() []model.Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return []model.Metric{*c.metric}
}

func (c *RandomMetricCollector) Update() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metric.value = c.rand.Float64()
}
