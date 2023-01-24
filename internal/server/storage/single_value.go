package storage

import (
	"context"

	"github.com/tony-spark/metrico/internal/model"
	"github.com/tony-spark/metrico/internal/server/models"
)

type SingleValueRepository struct {
	gauges   map[string]*models.GaugeValue
	counters map[string]*models.CounterValue
}

func NewSingleValueRepository() *SingleValueRepository {
	return &SingleValueRepository{
		gauges:   make(map[string]*models.GaugeValue),
		counters: make(map[string]*models.CounterValue),
	}
}

func (r SingleValueRepository) GetGaugeByName(_ context.Context, name string) (*models.GaugeValue, error) {
	return r.gauges[name], nil
}

func (r SingleValueRepository) SaveGauge(_ context.Context, name string, value float64) (*models.GaugeValue, error) {
	gauge, ok := r.gauges[name]
	if !ok {
		gauge = &models.GaugeValue{Name: name}
		r.gauges[name] = gauge
	}
	gauge.Value = value
	return gauge, nil
}

func (r SingleValueRepository) SaveAllGauges(ctx context.Context, gs []models.GaugeValue) error {
	for _, g := range gs {
		_, err := r.SaveGauge(ctx, g.Name, g.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r SingleValueRepository) GetCounterByName(_ context.Context, name string) (*models.CounterValue, error) {
	return r.counters[name], nil
}

func (r SingleValueRepository) AddAndSaveCounter(_ context.Context, name string, value int64) (*models.CounterValue, error) {
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

func (r SingleValueRepository) AddAndSaveAllCounters(ctx context.Context, cs []models.CounterValue) error {
	for _, c := range cs {
		_, err := r.AddAndSaveCounter(ctx, c.Name, c.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r SingleValueRepository) SaveCounter(_ context.Context, name string, value int64) (*models.CounterValue, error) {
	counter := &models.CounterValue{
		Name:  name,
		Value: value,
	}
	r.counters[name] = counter
	return counter, nil
}

func (r SingleValueRepository) GetAll(_ context.Context) ([]model.Metric, error) {
	ms := make([]model.Metric, 0, len(r.counters)+len(r.gauges))
	for _, c := range r.counters {
		ms = append(ms, *c)
	}
	for _, g := range r.gauges {
		ms = append(ms, *g)
	}
	return ms, nil
}
