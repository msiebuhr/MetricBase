package metrics

import (
	"fmt"
	"time"
)

// MetricValue are just the nameless components of metrics, as often
// used for storage and output.
type MetricValue struct {
	Value float64
	Time  time.Time
}

// Metric represents a metric by time, value and name.
type Metric struct {
	Name string
	MetricValue
}

// NewMetric generates a new Metric with the given paramters
func NewMetric(name string, value float64, when int64) *Metric {
	return &Metric{
		Name: name,
		MetricValue: MetricValue{
			Value: value,
			Time:  time.Unix(when, 0),
		},
	}
}

func (m Metric) GetMetricValue() MetricValue {
	return m.MetricValue
}

// Stringifies to match the Graphite protocol
func (m *Metric) String() string {
	return fmt.Sprintf("%s %v %d", m.Name, m.Value, m.Time.Unix())
}
