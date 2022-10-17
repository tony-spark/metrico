package models

type GaugeValue struct {
	Name  string
	Value float64
}

type CounterValue struct {
	Name  string
	Value int64
}
