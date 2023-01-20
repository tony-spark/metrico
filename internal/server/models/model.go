// Package models contains main structs and interfaces
package models

import (
	"fmt"

	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
)

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
	return model.GAUGE
}

func (g GaugeValue) Val() interface{} {
	return g.Value
}

func (g GaugeValue) String() string {
	return fmt.Sprint(g.Value)
}

func (c CounterValue) ID() string {
	return c.Name
}

func (c CounterValue) Type() string {
	return model.COUNTER
}

func (c CounterValue) Val() interface{} {
	return c.Value
}

func (c CounterValue) String() string {
	return fmt.Sprint(c.Value)
}

func FromDTO(mdto dto.Metric) model.Metric {
	switch mdto.MType {
	case model.GAUGE:
		return GaugeValue{
			Name:  mdto.ID,
			Value: *mdto.Value,
		}
	case model.COUNTER:
		return CounterValue{
			Name:  mdto.ID,
			Value: *mdto.Delta,
		}
	}
	return nil
}
