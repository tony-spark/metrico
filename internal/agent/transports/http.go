package transports

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/crypto"

	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
)

const (
	endpointSend          = "/update/{type}/{name}/{value}"
	endpointSendJSON      = "/update/"
	endpointSendJSONBatch = "/updates/"
)

type HTTPTransport struct {
	client    *resty.Client
	hasher    dto.Hasher
	encryptor crypto.Encryptor
}

type HTTPOption func(t *HTTPTransport)

func NewHTTP(baseURL string, options ...HTTPOption) Transport {
	var t HTTPTransport

	client := resty.New()
	client.SetBaseURL(baseURL)
	// TODO think about better timeout value
	client.SetTimeout(1 * time.Second)
	t.client = client

	for _, opt := range options {
		opt(&t)
	}

	return t
}

func WithHasher(h dto.Hasher) HTTPOption {
	return func(t *HTTPTransport) {
		t.hasher = h
	}
}

func WithEncryptor(e crypto.Encryptor) HTTPOption {
	return func(t *HTTPTransport) {
		t.encryptor = e
	}
}

func (h HTTPTransport) SendMetric(metric model.Metric) error {
	return h.sendJSON(metric)
}

func (h HTTPTransport) SendMetrics(mx []model.Metric) error {
	return h.sendJSONBatch(context.Background(), mx)
}

func (h HTTPTransport) SendMetricsWithContext(ctx context.Context, mx []model.Metric) error {
	return h.sendJSONBatch(ctx, mx)
}

func (h HTTPTransport) send(metric model.Metric) error {
	req := h.client.R().
		SetPathParam("type", metric.Type()).
		SetPathParam("name", metric.ID()).
		SetPathParam("value", metric.String()).
		SetHeader("Content-Type", "text/plain")
	resp, err := req.Post(endpointSend)
	if err != nil {
		return fmt.Errorf("could not send metric: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("send error: value not accepted %v response code: %v", req.URL, resp.StatusCode())
	}
	log.Info().Msgf("sent %v (%v) = %v", metric.ID(), metric.Type(), metric.String())
	return nil
}

func (h HTTPTransport) createDTO(metric model.Metric) (*dto.Metric, error) {
	d := dto.NewMetric(metric)
	if h.hasher != nil {
		var err error
		d.Hash, err = h.hasher.Hash(*d)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (h HTTPTransport) sendJSON(metric model.Metric) error {
	d, err := h.createDTO(metric)
	if err != nil {
		return err
	}
	req := h.client.R().
		SetBody(d)
	resp, err := req.Post(endpointSendJSON)
	if err != nil {
		return fmt.Errorf("could not send json: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("send error: value not accepted %v response code: %v", req.URL, resp.StatusCode())
	}
	log.Info().Msgf("sent %v (%v) = %v", metric.ID(), metric.Type(), metric.String())
	return nil
}

func (h HTTPTransport) sendJSONBatch(ctx context.Context, mx []model.Metric) error {
	var dtos []dto.Metric
	for _, m := range mx {
		mdto, err := h.createDTO(m)
		if err != nil {
			return err
		}
		dtos = append(dtos, *mdto)
	}
	req := h.client.R().
		SetContext(ctx)
	err := h.encodeInRequest(dtos, req)
	if err != nil {
		return err
	}

	resp, err := req.Post(endpointSendJSONBatch)
	if err != nil {
		return fmt.Errorf("could not send batch json: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("send error: metrics not accepted %v response code: %v", req.URL, resp.StatusCode())
	}
	for _, metric := range mx {
		log.Info().Msgf("sent in batch %v (%v) = %v", metric.ID(), metric.Type(), metric.String())
	}
	return nil
}

func (h HTTPTransport) encodeInRequest(obj interface{}, r *resty.Request) error {
	bs, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("could not marshal json: %w", err)
	}
	r.SetHeader("Content-Type", "application/json")
	if h.encryptor != nil {
		encrypted, err := h.encryptor.Encrypt(bs)
		if err != nil {
			return fmt.Errorf("could not encrypt message: %w", err)
		}
		r.SetHeader("X-Encrypted", "true")
		r.SetBody(encrypted)
	} else {
		r.SetBody(bs)
	}
	return nil
}
