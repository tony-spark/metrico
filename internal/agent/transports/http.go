package transports

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tony-spark/metrico/internal/agent/metrics"
	"log"
	"net/http"
	"time"
)

const (
	endpointSend = "/update/{type}/{name}/{value}"
)

type HTTPTransport struct {
	client *resty.Client
}

func NewHTTPTransport(baseURL string) *HTTPTransport {
	client := resty.New()
	client.SetBaseURL(baseURL)
	// TODO think about better timeout value
	client.SetTimeout(1 * time.Second)
	return &HTTPTransport{client}
}

func (h HTTPTransport) SendMetric(metric metrics.Metric) error {
	req := h.client.R().
		SetPathParam("type", metric.Type()).
		SetPathParam("name", metric.Name()).
		SetPathParam("value", metric.String()).
		SetHeader("Content-Type", "text/plain")
	resp, err := req.Post(endpointSend)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("send error: value not accepted %v response code: %v", req.URL, resp.StatusCode())
	}
	log.Printf("sent %v (%v) = %v\n", metric.Name(), metric.Type(), metric.String())
	return nil
}
