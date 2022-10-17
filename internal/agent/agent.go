package agent

import (
	"context"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"log"
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
	log.Println("poll")
	for _, collector := range a.collectors {
		collector.Update()
		for _, metric := range collector.Metrics() {
			log.Printf("got %v (%v) = %v\n", metric.Name(), metric.Type(), metric.String())
		}
	}
}

// TODO: if HTTP requests is taking too long, don't allow report goroutines to pile up
// TODO: what if report is running concurrently with poll?
func (a MetricsAgent) report() {
	log.Println("sending report...")
	for _, collector := range a.collectors {
		for _, metric := range collector.Metrics() {
			err := a.transport.SendMetric(metric)
			if err != nil {
				log.Println(err.Error())
				switch err.(type) {
				// TODO: move this logic to transport layer?
				case net.Error:
					log.Println("network error, interrupting current report...")
					return
				}
				continue
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
func (a MetricsAgent) Run(c context.Context) {
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
		case <-c.Done():
			// TODO: interrupt poll() and report()
			log.Println("Agent stopped via context")
			return
		}
	}
}
