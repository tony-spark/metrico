package models

import (
	"context"
	"io"
)

type GaugeRepository interface {
	GetGaugeByName(ctx context.Context, name string) (*GaugeValue, error)
	SaveGauge(ctx context.Context, name string, value float64) (*GaugeValue, error)
	SaveAllGauges(ctx context.Context, gs []GaugeValue) error
	GetAllGauges(ctx context.Context) ([]GaugeValue, error)
}

type CounterRepository interface {
	GetCounterByName(ctx context.Context, name string) (*CounterValue, error)
	AddAndSaveCounter(ctx context.Context, name string, value int64) (*CounterValue, error)
	AddAndSaveAllCounters(ctx context.Context, cs []CounterValue) error
	SaveCounter(ctx context.Context, name string, value int64) (*CounterValue, error)
	GetAllCounters(ctx context.Context) ([]CounterValue, error)
}

type DBManager interface {
	io.Closer
	Check(ctx context.Context) (bool, error)
	GaugeRepository() GaugeRepository
	CounterRepository() CounterRepository
}

type RepositoryPersistence interface {
	io.Closer
	Save(ctx context.Context, gr GaugeRepository, cr CounterRepository) error
	Load(ctx context.Context, gr GaugeRepository, cr CounterRepository) error
}
