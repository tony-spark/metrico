package model

import (
	"fmt"
)

const (
	COUNTER string = "counter"
	GAUGE   string = "gauge"
)

type Metric interface {
	fmt.Stringer
	ID() string
	Type() string
	Val() interface{}
}
