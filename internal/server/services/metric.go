package services

import (
	"context"
	"fmt"
	"github.com/tony-spark/metrico/internal/model"
	"github.com/tony-spark/metrico/internal/server/models"
)

type MetricService struct {
	gr         models.GaugeRepository
	cr         models.CounterRepository
	postUpdate func()
}

func NewMetricService(gr models.GaugeRepository, cr models.CounterRepository, postUpdate func()) *MetricService {
	return &MetricService{
		gr:         gr,
		cr:         cr,
		postUpdate: postUpdate,
	}
}

func (s MetricService) UpdateGauge(ctx context.Context, g models.GaugeValue) (gv *models.GaugeValue, err error) {
	gv, err = s.gr.Save(ctx, g.Name, g.Value)
	if err == nil && s.postUpdate != nil {
		s.postUpdate()
	}
	return
}

func (s MetricService) UpdateCounter(ctx context.Context, c models.CounterValue) (cv *models.CounterValue, err error) {
	cv, err = s.cr.AddAndSave(ctx, c.Name, c.Value)
	if err == nil && s.postUpdate != nil {
		s.postUpdate()
	}
	return
}

func (s MetricService) UpdateMetric(ctx context.Context, m model.Metric) (model.Metric, error) {
	switch m := m.(type) {
	case models.GaugeValue:
		return s.UpdateGauge(ctx, m)
	case models.CounterValue:
		return s.UpdateCounter(ctx, m)
	default:
		return nil, fmt.Errorf("unknown metric type")
	}
}

func (s MetricService) UpdateAll(ctx context.Context, gs []models.GaugeValue, cs []models.CounterValue) error {
	// TODO do we need single db transaction here?
	if len(gs) > 0 {
		err := s.gr.SaveAll(ctx, gs)
		if err != nil {
			return fmt.Errorf("could not save metris: %w", err)
		}
	}
	if len(cs) > 0 {
		err := s.cr.AddAndSaveAll(ctx, cs)
		if err != nil {
			return fmt.Errorf("could not save metris: %w", err)
		}
	}
	if s.postUpdate != nil {
		s.postUpdate()
	}
	return nil
}

func (s MetricService) Get(ctx context.Context, name string, mType string) (model.Metric, error) {
	switch mType {
	case model.GAUGE:
		g, err := s.gr.GetByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve gauge value: %w", err)
		}
		if g == nil {
			return nil, nil
		}
		return g, nil
	case model.COUNTER:
		c, err := s.cr.GetByName(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve counter value: %w", err)
		}
		if c == nil {
			return nil, nil
		}
		return c, nil
	default:
		return nil, fmt.Errorf("unknown metric type")
	}
}

func (s MetricService) GetAll(ctx context.Context) ([]model.Metric, error) {
	var ms []model.Metric
	gs, err := s.gr.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve all gauges: %w", err)
	}
	cs, err := s.cr.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve all counters: %w", err)
	}
	for _, g := range gs {
		ms = append(ms, g)
	}
	for _, c := range cs {
		ms = append(ms, c)
	}
	return ms, nil
}
