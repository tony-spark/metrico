package models

type GaugeRepository interface {
	GetByName(name string) (*GaugeValue, error)
	Save(name string, value float64) (*GaugeValue, error)
	GetAll() ([]*GaugeValue, error)
}

type CounterRepository interface {
	GetByName(name string) (*CounterValue, error)
	AddAndSave(name string, value int64) (*CounterValue, error)
	GetAll() ([]*CounterValue, error)
}
