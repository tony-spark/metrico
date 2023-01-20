// Package dto provides main interfaces and type for metric sending and receiving (e.g. DTOs)
package dto

import (
	"github.com/tony-spark/metrico/internal/model"
)

// Metric is a DTO with metric's data
type Metric struct {
	ID    string   `json:"id"`                          // metric's ID
	MType string   `json:"type" enums:"gauge,counter"`  // type of metric ("gauge" or "counter)
	Delta *int64   `json:"delta,omitempty"`             // value of counter metric
	Value *float64 `json:"value,omitempty"`             // value of gauge metric
	Hash  string   `json:"hash,omitempty" format:"HEX"` // object hash
}

// Hasher implementation is used to calculate and check DTO's hash
type Hasher interface {
	// Hash returns string (hex) representation of hash or error if hash can't be calculated for given Metric (e.g. inconsistent object)
	Hash(m Metric) (string, error)
	// Check returns hash check result or error if hash can't be checked for given Metric (e.g. inconsistent object)
	Check(m Metric) (bool, error)
}

// HasValue returns whether Metric is filled with value, depending on it's type
func (m Metric) HasValue() bool {
	switch m.MType {
	case model.GAUGE:
		return m.Value != nil
	case model.COUNTER:
		return m.Delta != nil
	}
	return false
}

// NewMetric creates a DTO from model object
func NewMetric(m model.Metric) *Metric {
	mdto := &Metric{
		ID:    m.ID(),
		MType: m.Type(),
		Delta: nil,
		Value: nil,
	}

	switch m.Type() {
	case model.GAUGE:
		var v float64
		v = m.Val().(float64)
		mdto.Value = &v
	case model.COUNTER:
		var d int64
		d = m.Val().(int64)
		mdto.Delta = &d
	}

	return mdto
}
