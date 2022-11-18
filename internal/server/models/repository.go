package models

import (
	"context"
	"github.com/tony-spark/metrico/internal/model"
	"io"
)

type MetricRepository interface {
	GetGaugeByName(ctx context.Context, name string) (*GaugeValue, error)
	SaveGauge(ctx context.Context, name string, value float64) (*GaugeValue, error)
	SaveAllGauges(ctx context.Context, gs []GaugeValue) error
	GetCounterByName(ctx context.Context, name string) (*CounterValue, error)
	AddAndSaveCounter(ctx context.Context, name string, value int64) (*CounterValue, error)
	AddAndSaveAllCounters(ctx context.Context, cs []CounterValue) error
	SaveCounter(ctx context.Context, name string, value int64) (*CounterValue, error)
	GetAll(ctx context.Context) ([]model.Metric, error)
}

type DBManager interface {
	io.Closer
	Check(ctx context.Context) (bool, error)
	MetricRepository() MetricRepository
}

type RepositoryPersistence interface {
	io.Closer
	Save(ctx context.Context, r MetricRepository) error
	Load(ctx context.Context, r MetricRepository) error
}
