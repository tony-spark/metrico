package storage

import (
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

func (r SingleValueGaugeRepository) GetByName(name string) (*models.GaugeValue, error) {
	return r.gauges[name], nil
}

func (r SingleValueGaugeRepository) GetAll() ([]*models.GaugeValue, error) {
	return maps.Values(r.gauges), nil
}

func (r SingleValueGaugeRepository) Save(name string, value float64) (*models.GaugeValue, error) {
	gauge, ok := r.gauges[name]
	if !ok {
		gauge = &models.GaugeValue{Name: name}
		r.gauges[name] = gauge
	}
	gauge.Value = value
	return gauge, nil
}

func (r SingleValueCounterRepository) GetByName(name string) (*models.CounterValue, error) {
	return r.counters[name], nil
}

func (r SingleValueCounterRepository) GetAll() ([]*models.CounterValue, error) {
	return maps.Values(r.counters), nil
}

func (r SingleValueCounterRepository) AddAndSave(name string, value int64) (*models.CounterValue, error) {
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
