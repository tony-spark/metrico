package models

import (
	"context"
	"io"
)

type GaugeRepository interface {
	GetByName(name string) (*GaugeValue, error)
	Save(name string, value float64) (*GaugeValue, error)
	GetAll() ([]*GaugeValue, error)
}

type CounterRepository interface {
	GetByName(name string) (*CounterValue, error)
	AddAndSave(name string, value int64) (*CounterValue, error)
	Save(name string, value int64) (*CounterValue, error)
	GetAll() ([]*CounterValue, error)
}

type DBRepository interface {
	io.Closer
	Check(ctx context.Context) (bool, error)
}

type RepositoryPersistence interface {
	io.Closer
	Save(gr GaugeRepository, cr CounterRepository) error
	Load(gr GaugeRepository, cr CounterRepository) error
}
