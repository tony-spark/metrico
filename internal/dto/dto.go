package dto

import "github.com/tony-spark/metrico/internal"

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
	case internal.GAUGE:
		return m.Value != nil
	case internal.COUNTER:
		return m.Delta != nil
	}
	return false
}
