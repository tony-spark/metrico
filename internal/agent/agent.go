package agent

import (
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"log"
	"net"
	"time"
)

var collectors []metrics.MetricCollector
var transp transports.Transport

func poll() {
	log.Println("poll")
	for _, collector := range collectors {
		collector.Update()
		for _, metric := range collector.Metrics() {
			log.Printf("got %v (%v) = %v\n", metric.Name(), metric.Type(), metric.String())
		}
	}
}

// TODO: if HTTP requests is taking too long, don't allow report goroutines to pile up
// TODO: what if report is running concurrently with poll?
func report() {
	log.Println("sending report...")
	for _, collector := range collectors {
		for _, metric := range collector.Metrics() {
			err := transp.SendMetric(metric)
			if err != nil {
				log.Println(err.Error())
				switch err.(type) {
				// TODO: move this logic to transport layer
				case net.Error:
					log.Println("network error, interrupting current report...")
					return
				}
				continue
			}
		}
	}
}

// Run runs agent for collecting metrics data and sending it to server
// TODO: way to stop agent (pass Context?)
func Run(pollInterval time.Duration, reportInterval time.Duration, transport transports.Transport) {
	collectors = append(collectors, metrics.NewMemoryMetricCollector(), metrics.NewRandomMetricCollector())
	transp = transport

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	defer func() {
		pollTicker.Stop()
		reportTicker.Stop()
	}()

	for {
		select {
		case <-pollTicker.C:
			go poll()
		case <-reportTicker.C:
			go report()
		}
	}
}
