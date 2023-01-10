package agent_test

import (
	"context"
	"time"

	"github.com/tony-spark/metrico/internal/agent"
	"github.com/tony-spark/metrico/internal/agent/metrics"
)

// This example runs agent with default poll interval, 15 seconds report interval, only with random value metric
func Example() {
	a := agent.New(
		agent.WithReportInterval(15*time.Second),
		agent.WithHTTPTransport("http://localhost:3000", ""),
		agent.WithCollectors([]metrics.MetricCollector{
			metrics.NewRandomMetricCollector(),
		}),
	)

	go a.Run(context.Background())

	time.Sleep(30 * time.Second)
}
