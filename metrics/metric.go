package metrics

import (
	"fmt"
)

// MetricValue are just the nameless components of metrics, as often
// used for storage and output.
type MetricValue struct {
	Time  int64
	Value float64
}

// Metric represents a metric by time, value and name.
type Metric struct {
	MetricValue
	Name string
}

// NewMetric generates a new Metric with the given paramters
func NewMetric(name string, value float64, time int64) *Metric {
	m := &Metric{
		Name: name,
	}
	m.Value = value
	m.Time = time
	return m
}

// Stringifies to match the Graphite protocol
func (m *Metric) String() string {
	return fmt.Sprintf("%s %v %d", m.Name, m.Value, m.Time)
}
