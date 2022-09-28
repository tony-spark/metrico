package models

type NamedValue struct {
	Name string
}

type GaugeValue struct {
	NamedValue
	Value float64
}

type CounterValue struct {
	NamedValue
	Value int64
}
