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
)

// MetricsAgent represents agent application
type MetricsAgent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	collectors     []metrics.MetricCollector
	transport      transports.Transport
	mu             *sync.Mutex
	cond           *sync.Cond
	sending        bool
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
		transport: transports.NewHTTP("http://127.0.0.1:8080"),
	}
	a.mu = new(sync.Mutex)
	a.cond = sync.NewCond(a.mu)

	for _, opt := range options {
		opt(&a)
	}

	return a
}

// WithTransport configures agent to use given transport to send metrics
func WithTransport(transport transports.Transport) Option {
	return func(a *MetricsAgent) {
		a.transport = transport
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
func WithCollectors(cs ...metrics.MetricCollector) Option {
	return func(a *MetricsAgent) {
		a.collectors = cs
	}
}

func (a MetricsAgent) poll() {
	log.Trace().Msg("poll")
	for _, collector := range a.collectors {
		collector.Update()
		for _, metric := range collector.Metrics() {
			log.Debug().Msgf("got %v (%v) = %v", metric.ID(), metric.Type(), metric.String())
		}
	}
}

func (a MetricsAgent) report() {
	log.Info().Msg("sending report")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), a.reportInterval)
	defer cancel()

	a.mu.Lock()
	defer a.mu.Unlock()

	a.sending = true

	var wg sync.WaitGroup

	for _, collector := range a.collectors {
		wg.Add(1)
		go func(c metrics.MetricCollector) {
			defer wg.Done()

			select {
			case <-timeoutCtx.Done():
				log.Warn().Msg("sending cancelled timeout")
				return
			default:
			}

			err := a.transport.SendMetricsWithContext(timeoutCtx, c.Metrics())
			if err != nil {
				log.Error().Err(err).Msg("could not send metrics")
				var ne net.Error
				if errors.As(err, &ne) {
					log.Info().Msg("network error, interrupting current report...")
					a.cond.Broadcast()
					return
				}
			}
		}(collector)
	}

	wg.Wait()
	a.sending = false
	a.cond.Broadcast()
}

// Run starts to collect metrics and send it via transport
//
// Note that Run blocks until given context is cancelled or Stop called
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
			go a.poll()
		case <-reportTicker.C:
			go a.report()
		case <-ctx.Done():
			log.Info().Msg("agent stopped via context")
			return
		}
	}
}

// Stop gracefully stops agent
func (a MetricsAgent) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.sending {
		a.cond.Wait()
	}
}
