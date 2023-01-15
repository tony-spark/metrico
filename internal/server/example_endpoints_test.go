package server_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func doRequest(method string, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error while configuring request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during request: %w", err)
	}
	return resp, nil
}

// Example_sendMetricJSON demonstrates how to send metric to server using JSON
func Example_sendMetricJSON() {
	json := `{
		"id" : "GaugeExample",
		"type" : "gauge",
		"value" : 1.23
	}`
	resp, err := doRequest(http.MethodPost, "http://localhost:8080/update/", []byte(json))
	if err != nil {
		log.Fatal().Err(err).Msg("error during request")
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

// Example_sendMetricJSON demonstrates how to send multiple metrics in batch
func Example_sendMetricsInBatch() {
	json := `[
		{
			"id" : "GaugeExample",
			"type" : "gauge",
			"value" : 1.5
		},
		{
			"id" : "CounterExample",
			"type" : "counter",
			"delta" : 4
		}
	]`
	resp, err := doRequest(http.MethodPost, "http://localhost:8080/updates/", []byte(json))
	if err != nil {
		log.Fatal().Err(err).Msg("error during request")
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

// Example_getMetric demonstrates how to get metric value
func Example_getMetric() {
	json := `{
		"id" : "GaugeExample",
		"type" : "gauge"
	}`
	resp, err := doRequest(http.MethodPost, "http://localhost:8080/updates/", []byte(json))
	if err != nil {
		log.Fatal().Err(err).Msg("error during request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read response body")
	}

	fmt.Println(string(body))
}
