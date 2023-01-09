package metrics

import (
	"time"

	"github.com/tony-spark/metrico/internal/model"
)

type DelayedCollectorProxy struct {
	delay time.Duration
	c     MetricCollector
}

func NewDelayedCollectorProxy(c MetricCollector, delay time.Duration) DelayedCollectorProxy {
	return DelayedCollectorProxy{
		c:     c,
		delay: delay,
	}
}

func (p DelayedCollectorProxy) Metrics() []model.Metric {
	return p.c.Metrics()
}

func (p DelayedCollectorProxy) Update() {
	time.Sleep(p.delay)
	p.c.Update()
}
