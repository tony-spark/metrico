package agent

import (
	"github.com/tony-spark/metrico/internal"
	"log"
	"time"
)

var collectors []internal.MetricCollector

func initCollectors() {
	collectors = append(collectors, NewMemoryMetricCollector(), NewRandomMetricCollector())
}

func poll() {
	log.Println("poll")
	for _, collector := range collectors {
		collector.Update()
		for _, metric := range collector.Metrics() {
			log.Println(metric.Name(), metric.String())
		}
	}
}

// TODO: if HTTP request is taking too long, don't allow report goroutines to pile up
// TODO: if report is running simultaneously with poll?
func report() {
	log.Println("report")
}

// Run runs agent for collecting mxs data and sending it to server
// TODO: way to stop agent (pass Context?)
func Run(pollInterval time.Duration, reportInterval time.Duration, serverAddress string) {
	initCollectors()
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
