package metrics

import "fmt"

type Metric interface {
	fmt.Stringer
	Name() string
	Type() string
}

type MetricCollector interface {
	Metrics() []Metric
	Update()
}
