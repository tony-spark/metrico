package grpc

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/rs/zerolog/log"
	pb "github.com/tony-spark/metrico/gen/pb/api"
	"github.com/tony-spark/metrico/internal/agent/transports"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Transport struct {
	client pb.MetricServiceClient
	hasher dto.Hasher
}

type Option func(t *Transport)

func WithHasher(h dto.Hasher) Option {
	return func(t *Transport) {
		t.hasher = h
	}
}

func NewTransport(addr string, opts ...Option) (transports.Transport, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("could not dial grpc: %w", err)
	}

	t := Transport{
		client: pb.NewMetricServiceClient(conn),
	}

	for _, opt := range opts {
		opt(&t)
	}

	return t, nil
}

func (t Transport) SendMetric(metric model.Metric) error {
	return nil
}

func (t Transport) SendMetrics(mx []model.Metric) error {
	return t.SendMetricsWithContext(context.Background(), mx)
}

func (t Transport) SendMetricsWithContext(ctx context.Context, mx []model.Metric) error {
	uc, err := t.client.Update(ctx)
	if err != nil {
		return fmt.Errorf("could not init grpc stream: %w", err)
	}
	for _, m := range mx {
		var mt *pb.Metric
		mt, err = t.createDTO(m)
		if err != nil {
			log.Error().Err(err).Msg("could not create dto")
			continue
		}
		err = uc.Send(mt)
		if err != nil {
			log.Error().Err(err).Msg("could not send")
			continue
		}
		log.Info().Msgf("sent  %v (%v) = %v", m.ID(), m.Type(), m.String())
	}
	var r *pb.Response
	r, err = uc.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("could not close stream: %w", err)
	}
	if r.Status == pb.Status_ERROR {
		return fmt.Errorf("error sending metrics: %s", r.GetError())
	}
	return nil
}

func (t Transport) createDTO(metric model.Metric) (*pb.Metric, error) {
	d := dto.NewMetric(metric)
	var hash []byte
	if t.hasher != nil {
		var err error
		d.Hash, err = t.hasher.Hash(*d)
		if err != nil {
			return nil, err
		}
		hash, err = hex.DecodeString(d.Hash)
		if err != nil {
			return nil, err
		}
	}
	var mt pb.MetricType
	switch metric.Type() {
	case model.GAUGE:
		mt = pb.MetricType_GAUGE
	case model.COUNTER:
		mt = pb.MetricType_COUNTER
	}
	return &pb.Metric{
		Id:    d.ID,
		Type:  mt,
		Delta: d.Delta,
		Value: d.Value,
		Hash:  hash,
	}, nil
}
