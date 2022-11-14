package transports

import "github.com/tony-spark/metrico/internal/model"

type Dummy struct {
}

func NewDummy() Dummy {
	return Dummy{}
}

func (d Dummy) SendMetric(_ model.Metric) error {
	return nil
}

func (d Dummy) SendMetrics(_ []model.Metric) error {
	return nil
}
