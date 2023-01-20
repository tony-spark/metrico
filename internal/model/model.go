// Package model provides core interfaces and constants
package model

import (
	"fmt"
)

// Metric types
const (
	COUNTER string = "counter" // counter metric has int64 value
	GAUGE   string = "gauge"   // gauge metric has float64 value
)

// Metric is a main model interface
type Metric interface {
	fmt.Stringer
	// ID returns metric's id
	ID() string
	// Type returns metric's type (one of the constants)
	Type() string
	// Val returns metric's value
	Val() interface{}
}
