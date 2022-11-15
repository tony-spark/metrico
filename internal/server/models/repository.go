package models

import (
	"context"
	"io"
)

type GaugeRepository interface {
	GetByName(ctx context.Context, name string) (*GaugeValue, error)
	Save(ctx context.Context, name string, value float64) (*GaugeValue, error)
	SaveAll(ctx context.Context, gs []GaugeValue) error
	GetAll(ctx context.Context) ([]GaugeValue, error)
}

type CounterRepository interface {
	GetByName(ctx context.Context, name string) (*CounterValue, error)
	AddAndSave(ctx context.Context, name string, value int64) (*CounterValue, error)
	AddAndSaveAll(ctx context.Context, cs []CounterValue) error
	Save(ctx context.Context, name string, value int64) (*CounterValue, error)
	GetAll(ctx context.Context) ([]CounterValue, error)
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
