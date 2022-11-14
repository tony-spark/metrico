package agent

import (
	"context"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"testing"
	"time"
)

func TestAgent(t *testing.T) {
	a := NewMetricsAgent(
		1*time.Second,
		5*time.Second,
		transports.NewDummy(),
		[]metrics.MetricCollector{
			metrics.NewDelayedCollectorProxy(metrics.NewRandomMetricCollector(), 1*time.Second),
		})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go a.Run(ctx)

	time.Sleep(15 * time.Second)
}
