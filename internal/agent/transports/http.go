package transports

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
	endpointSend = "/update/{type}/{name}/{value}"
)

type HttpTransport struct {
	client *resty.Client
}

func NewHttpTransport(baseURL string) *HttpTransport {
	client := resty.New()
	client.SetBaseURL(baseURL)
	// TODO think about better timeout value
	client.SetTimeout(1 * time.Second)
	return &HttpTransport{client}
}

func (h HttpTransport) SendMetric(metric metrics.Metric) error {
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
		err = errors.New(fmt.Sprintf("send error: value not accepted %v response code: %v", req.URL, resp.StatusCode()))
		return err
	}
	log.Printf("sent %v (%v) = %v\n", metric.Name(), metric.Type(), metric.String())
	return nil
}
