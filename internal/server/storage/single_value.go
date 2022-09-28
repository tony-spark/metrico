package storage

import "github.com/tony-spark/metrico/internal/server/models"

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

func (r SingleValueGaugeRepository) Save(name string, value float64) (*models.GaugeValue, error) {
	gauge, ok := r.gauges[name]
	if !ok {
		gauge = &models.GaugeValue{NamedValue: models.NamedValue{Name: name}}
	}
	gauge.Value = value
	return gauge, nil
}

func (r SingleValueCounterRepository) GetByName(name string) (*models.CounterValue, error) {
	return r.counters[name], nil
}

func (r SingleValueCounterRepository) AddAndSave(name string, value int64) (*models.CounterValue, error) {
	counter, ok := r.counters[name]
	if !ok {
		counter = &models.CounterValue{
			NamedValue: models.NamedValue{Name: name},
			Value:      0,
		}
	}
	counter.Value += value
	return counter, nil
}
