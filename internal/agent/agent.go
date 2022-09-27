package agent

import (
	"github.com/tony-spark/metrico/internal"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var collectors []internal.MetricCollector
var client http.Client
var address string

func poll() {
	log.Println("poll")
	for _, collector := range collectors {
		collector.Update()
		for _, metric := range collector.Metrics() {
			log.Println(metric.Name(), metric.String())
		}
	}
}

// TODO: if HTTP requests is taking too long, don't allow report goroutines to pile up
// TODO: if report is running simultaneously with poll?
// TODO: detect server shutdown
func report() {
	log.Println("report")
	for _, collector := range collectors {
		collector.Update()
		for _, metric := range collector.Metrics() {
			err := sendMetric(metric)
			if err != nil {
				continue
			}
		}
	}
}

func sendMetric(metric internal.Metric) error {
	endpoint := address + "/update/" + metric.Type() + "/" + metric.Name() + "/" + metric.String()
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(""))
	if err != nil {
		log.Println(endpoint, err)
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(req.RequestURI, err)
		return err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading (empty) body", err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(req.RequestURI, resp.StatusCode)
	}
	return nil
}

// Run runs agent for collecting mxs data and sending it to server
// TODO: way to stop agent (pass Context?)
func Run(pollInterval time.Duration, reportInterval time.Duration, serverAddress string) {
	address = serverAddress
	collectors = append(collectors, NewMemoryMetricCollector(), NewRandomMetricCollector())
	client = http.Client{}
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
