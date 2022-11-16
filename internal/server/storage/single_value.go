package storage

import (
	"context"
	"github.com/tony-spark/metrico/internal/server/models"
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

func (r SingleValueGaugeRepository) GetGaugeByName(_ context.Context, name string) (*models.GaugeValue, error) {
	return r.gauges[name], nil
}

func (r SingleValueGaugeRepository) GetAllGauges(_ context.Context) ([]models.GaugeValue, error) {
	gs := make([]models.GaugeValue, 0, len(r.gauges))
	for _, g := range r.gauges {
		gs = append(gs, *g)
	}
	return gs, nil
}

func (r SingleValueGaugeRepository) SaveGauge(_ context.Context, name string, value float64) (*models.GaugeValue, error) {
	gauge, ok := r.gauges[name]
	if !ok {
		gauge = &models.GaugeValue{Name: name}
		r.gauges[name] = gauge
	}
	gauge.Value = value
	return gauge, nil
}

func (r SingleValueGaugeRepository) SaveAllGauges(ctx context.Context, gs []models.GaugeValue) error {
	for _, g := range gs {
		r.SaveGauge(ctx, g.Name, g.Value)
	}
	return nil
}

func (r SingleValueCounterRepository) GetCounterByName(_ context.Context, name string) (*models.CounterValue, error) {
	return r.counters[name], nil
}

func (r SingleValueCounterRepository) GetAllCounters(_ context.Context) ([]models.CounterValue, error) {
	cs := make([]models.CounterValue, 0, len(r.counters))
	for _, c := range r.counters {
		cs = append(cs, *c)
	}
	return cs, nil
}

func (r SingleValueCounterRepository) AddAndSaveCounter(_ context.Context, name string, value int64) (*models.CounterValue, error) {
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

func (r SingleValueCounterRepository) AddAndSaveAllCounters(ctx context.Context, cs []models.CounterValue) error {
	for _, c := range cs {
		r.AddAndSaveCounter(ctx, c.Name, c.Value)
	}
	return nil
}

func (r SingleValueCounterRepository) SaveCounter(_ context.Context, name string, value int64) (*models.CounterValue, error) {
	counter := &models.CounterValue{
		Name:  name,
		Value: value,
	}
	r.counters[name] = counter
	return counter, nil
}
