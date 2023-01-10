// Package agent contains implementation of metrics agent - application to collect metrics and send it to server
package agent

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"github.com/tony-spark/metrico/internal/hash"
)

// MetricsAgent represents agent application
type MetricsAgent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	collectors     []metrics.MetricCollector
	transport      transports.Transport
}

// Option represents option function for agent configuration
type Option func(a *MetricsAgent)

// New creates agent with given options
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

// WithHTTPTransport configures agent to send metrics to given URL via HTTP.
// If hashKey is not empty, hash will be calculated during sending metrics
func WithHTTPTransport(url string, hashKey string) Option {
	return func(a *MetricsAgent) {
		if len(hashKey) > 0 {
			a.transport = transports.NewHTTPTransportHashed(url, hash.NewSha256Hmac(hashKey))
		} else {
			a.transport = transports.NewHTTPTransport(url)
		}
	}
}

// WithPollInterval configures agent to update metrics at a given interval
func WithPollInterval(interval time.Duration) Option {
	return func(a *MetricsAgent) {
		a.pollInterval = interval
	}
}

// WithReportInterval configures agent to send metrics at a given interval
func WithReportInterval(interval time.Duration) Option {
	return func(a *MetricsAgent) {
		a.reportInterval = interval
	}
}

// WithCollectors configures agent with given set of metrics collectors
func WithCollectors(cs []metrics.MetricCollector) Option {
	return func(a *MetricsAgent) {
		a.collectors = cs
	}
}

// NewMetricsAgent creates new agent with given pollInterval, reportInterval, transport and collectors
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

func (a MetricsAgent) report(ctx context.Context) {
	log.Info().Msg("sending report")
	timeoutCtx, cancel := context.WithTimeout(ctx, a.reportInterval)
	defer cancel()

	var wg sync.WaitGroup

	for _, collector := range a.collectors {
		wg.Add(1)
		go func(c metrics.MetricCollector) {
			defer wg.Done()

			select {
			case <-timeoutCtx.Done():
				log.Warn().Msg("sending cancelled via context (timeout?)")
				return
			default:
			}

			err := a.transport.SendMetricsWithContext(timeoutCtx, c.Metrics())
			if err != nil {
				log.Error().Err(err).Msg("could not send metrics")
				var ne net.Error
				if errors.As(err, &ne) {
					log.Info().Msg("network error, interrupting current report...")
					return
				}
			}
		}(collector)
	}

	wg.Wait()
}

// Run starts to collect metrics and send it via transport
//
// Note that Run blocks until given context is done
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
