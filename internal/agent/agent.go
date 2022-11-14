package agent

import (
	"context"
	"errors"
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

func (a MetricsAgent) poll(ctx context.Context) {
	log.Trace().Msg("poll")
	for _, collector := range a.collectors {
		select {
		case <-ctx.Done():
			log.Warn().Msg("poll cancelled via context")
			return
		default:
		}

		collector.Update()
		for _, metric := range collector.Metrics() {
			log.Debug().Msgf("got %v (%v) = %v", metric.ID(), metric.Type(), metric.String())
		}
	}
}

// TODO: if HTTP requests is taking too long, don't allow report goroutines to pile up
// TODO: what if report is running concurrently with poll?
func (a MetricsAgent) report(ctx context.Context) {
	log.Info().Msg("sending report")
	for _, collector := range a.collectors {
		select {
		case <-ctx.Done():
			log.Warn().Msg("sending cancelled via context")
			return
		default:
		}

		err := a.transport.SendMetrics(collector.Metrics())
		if err != nil {
			log.Error().Err(err).Msg("could not send metrics")
			var ne net.Error
			if errors.As(err, &ne) {
				log.Info().Msg("network error, interrupting current report...")
				return
			}
		}
	}
}

func NewMetricsAgent(pollInterval time.Duration, reportInterval time.Duration, transport transports.Transport, collectors []metrics.MetricCollector) *MetricsAgent {
	return &MetricsAgent{
		transport:      transport,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		collectors:     collectors,
	}
}

// Run starts collecting metrics and sending it via transport
func (a MetricsAgent) Run(ctx context.Context) {
	pollTicker := time.NewTicker(a.pollInterval)
	reportTicker := time.NewTicker(a.reportInterval)
	defer func() {
		pollTicker.Stop()
		reportTicker.Stop()
	}()

	for {
		select {
		case <-pollTicker.C:
			go a.poll(ctx)
		case <-reportTicker.C:
			go a.report(ctx)
		case <-ctx.Done():
			log.Info().Msg("Agent stopped via context")
			return
		}
	}
}
