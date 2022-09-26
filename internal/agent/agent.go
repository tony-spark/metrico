package agent

import (
	"github.com/tony-spark/metrico/internal"
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
			req, err := newRequestWithMetric(metric)
			if err != nil {
				continue
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(req.RequestURI, err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Println(req.RequestURI, resp.StatusCode)
			}
		}
	}
}

func newRequestWithMetric(metric internal.Metric) (request *http.Request, err error) {
	endpoint := address + "/update/" + metric.Type() + "/" + metric.Name() + "/" + metric.String()
	request, err = http.NewRequest(http.MethodPost, endpoint, strings.NewReader(""))
	if err != nil {
		log.Println(endpoint, err)
		return nil, err
	}
	request.Header.Add("Content-Type", "text/plain")
	return
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
