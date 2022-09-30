package agent

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"log"
	"net/http"
	"time"
)

const (
	endpoint = "/update/{type}/{name}/{value}"
)

var collectors []metrics.MetricCollector
var client *resty.Client

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
// TODO: detect server shutdown (interrupt current report() and try next time)
func report() {
	log.Println("sending report...")
	for _, collector := range collectors {
		for _, metric := range collector.Metrics() {
			err := sendMetric(metric)
			if err != nil {
				log.Println(err.Error())
				continue
			}
		}
	}
}

func sendMetric(metric metrics.Metric) error {
	req := client.R().
		SetPathParam("type", metric.Type()).
		SetPathParam("name", metric.Name()).
		SetPathParam("value", metric.String()).
		SetHeader("Content-Type", "text/plain")
	resp, err := req.Post(endpoint)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		err = errors.New(fmt.Sprintf("send error: value not accepted %v response code: %v", req.URL, resp.StatusCode()))
		return err
	}
	log.Printf("sent %v (%v) = %v\n", metric.Name(), metric.Type(), metric.String())
	return nil
}

// Run runs agent for collecting mxs data and sending it to server
// TODO: way to stop agent (pass Context?)
func Run(pollInterval time.Duration, reportInterval time.Duration, baseURL string) {
	collectors = append(collectors, metrics.NewMemoryMetricCollector(), metrics.NewRandomMetricCollector())

	client = resty.New()
	client.SetBaseURL(baseURL)
	client.SetTimeout(1 * time.Second)

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
