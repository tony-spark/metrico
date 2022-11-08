package models

import (
	"github.com/tony-spark/metrico/internal"
	"github.com/tony-spark/metrico/internal/dto"
)

type Metric interface {
	ID() string
	Type() string
	V() interface{}
	ToDTO() *dto.Metric
}

type GaugeValue struct {
	Name  string
	Value float64
}

type CounterValue struct {
	Name  string
	Value int64
}

func (g GaugeValue) ID() string {
	return g.Name
}

func (g GaugeValue) Type() string {
	return internal.GAUGE
}

func (g GaugeValue) V() interface{} {
	return g.Value
}

func (g GaugeValue) ToDTO() *dto.Metric {
	v := g.Value
	return &dto.Metric{
		ID:    g.Name,
		MType: internal.GAUGE,
		Value: &v,
	}
}

func (c CounterValue) ID() string {
	return c.Name
}

func (c CounterValue) Type() string {
	return internal.COUNTER
}

func (c CounterValue) V() interface{} {
	return c.Value
}

func (c CounterValue) ToDTO() *dto.Metric {
	d := c.Value
	return &dto.Metric{
		ID:    c.Name,
		MType: internal.COUNTER,
		Delta: &d,
	}
}

func FromDTO(mdto dto.Metric) Metric {
	switch mdto.MType {
	case internal.GAUGE:
		return GaugeValue{
			Name:  mdto.ID,
			Value: *mdto.Value,
		}
	case internal.COUNTER:
		return CounterValue{
			Name:  mdto.ID,
			Value: *mdto.Delta,
		}
	}
	return nil
}
