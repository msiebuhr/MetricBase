package MetricBase

import (
	"fmt"
)

// Data structures
type MetricValues struct {
	Time  int64
	Value float64
}

type Metric struct {
	MetricValues
	Name string
}

// Requests
type DataRequest struct {
	Name   string
	From   int64
	To     int64
	Result chan MetricValues
}

// Interfaces
type Backend interface {
	Start()
	Stop()
	AddMetrics(chan Metric)
	GetRawData(string, int64, int64, chan MetricValues)
	GetMetricsList(chan string)
}

type Frontend interface {
	SetBackend(Backend)
	Start()
	Stop()
}

// Metric helper functions
func NewMetric(name string, value float64, time int64) *Metric {
	m := &Metric{
		Name: name,
	}
	m.Value = value
	m.Time = time
	return m
}

// Plain stringification
func (m *Metric) String() string {
	return fmt.Sprintf("%s %v %d", m.Name, m.Value, m.Time)
}
