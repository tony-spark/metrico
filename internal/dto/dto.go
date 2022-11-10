package dto

import (
	"github.com/tony-spark/metrico/internal/model"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type Hasher interface {
	Hash(m Metric) (string, error)
	Check(m Metric) (bool, error)
}

func (m Metric) HasValue() bool {
	switch m.MType {
	case model.GAUGE:
		return m.Value != nil
	case model.COUNTER:
		return m.Delta != nil
	}
	return false
}

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
