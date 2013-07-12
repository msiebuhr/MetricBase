package MetricBase

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
type AddRequest struct {
	Data chan Metric
}

type ListRequest struct {
	Result chan string
}

type DataRequest struct {
	Name   string
	From   int64
	To     int64
	Result chan MetricValues
}

// TODO: Helpers to generate these requests

// Interfaces
type Backend interface {
	SetBackend(Backend)
	Start()
	Stop()
	Add(AddRequest)
	List(ListRequest)
	Data(DataRequest)
}

type Frontend interface {
	SetBackend(Backend)
	Start()
	Stop()
}

// Metric helper functions

// Plain stringification
func (m *Metric) String() string {
	return fmt.Sprintf("%s %v %d", m.Name, m.Value, m.Time)
}

// Parse a line of graphite text format and return a new Metric
func ParseGraphiteLine(raw string) (error, Metric) {
	// Find newline-rune
	fields := strings.Fields(raw)

	newMetric := Metric{}

	// A line must at least contain <metric.name> <timestamp> <value> <tag=value>+
	if len(fields) != 3 {
		return errors.New("Invalid line"), newMetric
	}

	// Convert name
	newMetric.Name = string(fields[0])

	// Parse out value
	value, err := strconv.ParseFloat(string(fields[1]), 64)
	if err != nil {
		return err, newMetric
	}
	newMetric.Value = value

	// Parse out timestamp
	time, err := strconv.ParseInt(string(fields[2]), 10, 64)
	if err != nil {
		return err, newMetric
	}
	newMetric.Time = time

	return nil, newMetric
}
