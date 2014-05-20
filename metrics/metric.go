package metrics

import (
	"fmt"
)

// MetricValue are just the nameless components of metrics, as often
// used for storage and output.
type MetricValue struct {
	Value float64
	Time  int64
}

// Metric represents a metric by time, value and name.
type Metric struct {
	Name string
	MetricValue
}

// NewMetric generates a new Metric with the given paramters
func NewMetric(name string, value float64, time int64) *Metric {
	return &Metric{
		Name: name,
		MetricValue: MetricValue{
			Value: value,
			Time:  time,
		},
	}
}

// Stringifies to match the Graphite protocol
func (m *Metric) String() string {
	return fmt.Sprintf("%s %v %d", m.Name, m.Value, m.Time)
}
