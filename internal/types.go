package internal

import "fmt"

const (
	COUNTER = "counter"
	GAUGE   = "gauge"
)

type Metric interface {
	fmt.Stringer
	Name() string
	Type() string
}

type MetricCollector interface {
	Metrics() []Metric
	Update()
}
