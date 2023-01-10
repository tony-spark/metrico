package agent_test

import (
	"context"
	"testing"
	"time"

	"github.com/tony-spark/metrico/internal/agent"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/agent/transports"
)

func TestAgentRace(t *testing.T) {
	a := agent.NewMetricsAgent(
		1*time.Second,
		2*time.Second,
		transports.NewDelayed(transports.NewDummy(), 2*time.Second),
		[]metrics.MetricCollector{
			metrics.NewDelayedCollectorProxy(metrics.NewMemoryMetricCollector(), 1*time.Second),
			metrics.NewDelayedCollectorProxy(metrics.NewRandomMetricCollector(), 1*time.Second),
			metrics.NewDelayedCollectorProxy(metrics.NewPsUtilMetricsCollector(), 1*time.Second),
		})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go a.Run(ctx)

	time.Sleep(12 * time.Second)
}
