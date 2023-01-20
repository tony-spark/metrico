package transports

import (
	"context"
	"time"

	"github.com/tony-spark/metrico/internal/model"
)

type Dummy struct {
}

type Delayed struct {
	t     Transport
	delay time.Duration
}

func NewDelayed(t Transport, delay time.Duration) Delayed {
	return Delayed{
		t:     t,
		delay: delay,
	}
}

func (d Delayed) SendMetric(metric model.Metric) error {
	time.Sleep(d.delay)
	return d.t.SendMetric(metric)
}

func (d Delayed) SendMetrics(mx []model.Metric) error {
	return d.SendMetricsWithContext(context.Background(), mx)
}

func (d Delayed) SendMetricsWithContext(_ context.Context, mx []model.Metric) error {
	time.Sleep(d.delay)
	return d.t.SendMetrics(mx)
}

func NewDummy() Dummy {
	return Dummy{}
}

func (d Dummy) SendMetric(_ model.Metric) error {
	return nil
}

func (d Dummy) SendMetrics(_ []model.Metric) error {
	return nil
}

func (d Dummy) SendMetricsWithContext(_ context.Context, _ []model.Metric) error {
	return nil
}
