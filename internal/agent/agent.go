package agent

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"net"
	"time"
)

type MetricsAgent struct {
	collectors     []metrics.MetricCollector
	transport      transports.Transport
	pollInterval   time.Duration
	reportInterval time.Duration
}

func (a MetricsAgent) poll() {
	log.Trace().Msg("poll")
	for _, collector := range a.collectors {
		collector.Update()
		for _, metric := range collector.Metrics() {
			log.Debug().Msgf("got %v (%v) = %v", metric.Name(), metric.Type(), metric.String())
		}
	}
}

// TODO: if HTTP requests is taking too long, don't allow report goroutines to pile up
// TODO: what if report is running concurrently with poll?
func (a MetricsAgent) report() {
	log.Info().Msg("sending report")
	for _, collector := range a.collectors {
		err := a.transport.SendMetrics(collector.Metrics())
		if err != nil {
			log.Error().Msg(err.Error())
			switch err.(type) {
			// TODO: move this logic to transport layer?
			case net.Error:
				log.Info().Msg("network error, interrupting current report...")
				return
			}
		}
	}
}

func NewMetricsAgent(pollInterval time.Duration, reportInterval time.Duration, transport transports.Transport) *MetricsAgent {
	return &MetricsAgent{
		transport:      transport,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
}

// Run starts collecting metrics and sending it via transport
func (a MetricsAgent) Run(ctx context.Context) {
	a.collectors = append(a.collectors, metrics.NewMemoryMetricCollector(), metrics.NewRandomMetricCollector())

	pollTicker := time.NewTicker(a.pollInterval)
	reportTicker := time.NewTicker(a.reportInterval)
	defer func() {
		pollTicker.Stop()
		reportTicker.Stop()
	}()

	for {
		select {
		case <-pollTicker.C:
			go a.poll()
		case <-reportTicker.C:
			go a.report()
		case <-ctx.Done():
			// TODO: interrupt poll() and report()
			log.Info().Msg("Agent stopped via context")
			return
		}
	}
}
