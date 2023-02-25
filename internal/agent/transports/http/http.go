package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"github.com/tony-spark/metrico/internal/crypto"

	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
)

const (
	endpointSend          = "/update/{type}/{name}/{value}"
	endpointSendJSON      = "/update/"
	endpointSendJSONBatch = "/updates/"
)

type Transport struct {
	client    *resty.Client
	hasher    dto.Hasher
	encryptor crypto.Encryptor
	clientIP  string
}

type Option func(t *Transport)

func NewTransport(baseURL string, options ...Option) transports.Transport {
	var t Transport

	client := resty.New()
	client.SetBaseURL(baseURL)
	// TODO think about better timeout value
	client.SetTimeout(1 * time.Second)
	t.client = client

	for _, opt := range options {
		opt(&t)
	}

	t.clientIP = getClientIP(baseURL)

	return t
}

func WithHasher(h dto.Hasher) Option {
	return func(t *Transport) {
		t.hasher = h
	}
}

func WithEncryptor(e crypto.Encryptor) Option {
	return func(t *Transport) {
		t.encryptor = e
	}
}

func (h Transport) SendMetric(metric model.Metric) error {
	return h.sendJSON(metric)
}

func (h Transport) SendMetrics(mx []model.Metric) error {
	return h.sendJSONBatch(context.Background(), mx)
}

func (h Transport) SendMetricsWithContext(ctx context.Context, mx []model.Metric) error {
	return h.sendJSONBatch(ctx, mx)
}

func (h Transport) send(metric model.Metric) error {
	req := h.client.R().
		SetPathParam("type", metric.Type()).
		SetPathParam("name", metric.ID()).
		SetPathParam("value", metric.String()).
		SetHeader("Content-Type", "text/plain")

	req.SetHeader("X-Real-IP", h.clientIP)
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

func (h Transport) createDTO(metric model.Metric) (*dto.Metric, error) {
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

func (h Transport) sendJSON(metric model.Metric) error {
	d, err := h.createDTO(metric)
	if err != nil {
		return err
	}
	req := h.client.R().
		SetBody(d)
	req.SetHeader("X-Real-IP", h.clientIP)
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

func (h Transport) sendJSONBatch(ctx context.Context, mx []model.Metric) error {
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

	req.SetHeader("X-Real-IP", h.clientIP)
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

func (h Transport) encodeInRequest(obj interface{}, r *resty.Request) error {
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

func getClientIP(URL string) string {
	u, err := url.Parse(URL)
	if err != nil {
		log.Error().Err(err).Msg("could not parse URL")
	}
	hostname := strings.TrimPrefix(u.Hostname(), "www.")
	conn, err := net.Dial("udp", hostname+":80")
	if err != nil {
		log.Error().Err(err).Msg("could not dial address to discover own IP")
	}
	defer func() {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				log.Error().Err(err).Msg("error closing connection")
			}
		}
	}()

	if conn == nil {
		return ""
	}

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
