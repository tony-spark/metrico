package storage

import (
	"context"
	"github.com/tony-spark/metrico/internal/server/models"
	"golang.org/x/exp/maps"
)

type SingleValueGaugeRepository struct {
	gauges map[string]*models.GaugeValue
}

type SingleValueCounterRepository struct {
	counters map[string]*models.CounterValue
}

func NewSingleValueGaugeRepository() *SingleValueGaugeRepository {
	return &SingleValueGaugeRepository{
		gauges: make(map[string]*models.GaugeValue),
	}
}

func NewSingleValueCounterRepository() *SingleValueCounterRepository {
	return &SingleValueCounterRepository{
		counters: make(map[string]*models.CounterValue),
	}
}

func (r SingleValueGaugeRepository) GetByName(ctx context.Context, name string) (*models.GaugeValue, error) {
	return r.gauges[name], nil
}

func (r SingleValueGaugeRepository) GetAll(ctx context.Context) ([]*models.GaugeValue, error) {
	return maps.Values(r.gauges), nil
}

func (r SingleValueGaugeRepository) Save(ctx context.Context, name string, value float64) (*models.GaugeValue, error) {
	gauge, ok := r.gauges[name]
	if !ok {
		gauge = &models.GaugeValue{Name: name}
		r.gauges[name] = gauge
	}
	gauge.Value = value
	return gauge, nil
}

func (r SingleValueCounterRepository) GetByName(ctx context.Context, name string) (*models.CounterValue, error) {
	return r.counters[name], nil
}

func (r SingleValueCounterRepository) GetAll(ctx context.Context) ([]*models.CounterValue, error) {
	return maps.Values(r.counters), nil
}

func (r SingleValueCounterRepository) AddAndSave(ctx context.Context, name string, value int64) (*models.CounterValue, error) {
	counter, ok := r.counters[name]
	if !ok {
		counter = &models.CounterValue{
			Name:  name,
			Value: 0,
		}
		r.counters[name] = counter
	}
	counter.Value += value
	return counter, nil
}

func (r SingleValueCounterRepository) Save(ctx context.Context, name string, value int64) (*models.CounterValue, error) {
	counter := &models.CounterValue{
		Name:  name,
		Value: value,
	}
	r.counters[name] = counter
	return counter, nil
}
