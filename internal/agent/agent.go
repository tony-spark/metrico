package agent

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"github.com/tony-spark/metrico/internal/hash"
	"net"
	"time"
)

type MetricsAgent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	collectors     []metrics.MetricCollector
	transport      transports.Transport
}

type Option func(a *MetricsAgent)

func New(options ...Option) MetricsAgent {
	a := MetricsAgent{
		pollInterval:   2 * time.Second,
		reportInterval: 10 * time.Second,
		collectors: []metrics.MetricCollector{
			metrics.NewMemoryMetricCollector(),
			metrics.NewRandomMetricCollector(),
			metrics.NewPsUtilMetricsCollector(),
		},
		transport: transports.NewHTTPTransport("http://127.0.0.1:8080"),
	}

	for _, opt := range options {
		opt(&a)
	}

	return a
}

func WithHTTPTransport(url string, hashKey string) Option {
	return func(a *MetricsAgent) {
		if len(hashKey) > 0 {
			a.transport = transports.NewHTTPTransportHashed(url, hash.NewSha256Hmac(hashKey))
		} else {
			a.transport = transports.NewHTTPTransport(url)
		}
	}
}

func WithPollInterval(interval time.Duration) Option {
	return func(a *MetricsAgent) {
		a.pollInterval = interval
	}
}

func WithReportInterval(interval time.Duration) Option {
	return func(a *MetricsAgent) {
		a.reportInterval = interval
	}
}

func WithCollectors(cs []metrics.MetricCollector) Option {
	return func(a *MetricsAgent) {
		a.collectors = cs
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
