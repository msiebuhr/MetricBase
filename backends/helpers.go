package backends

import (
	"github.com/msiebuhr/MetricBase"
)

// GetMetricsAsList fetches all the metrics in a backend and returns it as a
// list.
func GetMetricsAsList(backend MetricBase.Backend) []string {
	list := make([]string, 0)
	out := make(chan string)

	backend.GetMetricsList(out)
	for name := range out {
		list = append(list, name)
	}
	return list
}

// GetDataAsList fetches the relevant data and returns it as a list.
func GetDataAsList(backend MetricBase.Backend, name string, from, to int64) []MetricBase.MetricValues {
	list := make([]MetricBase.MetricValues, 0)
	out := make(chan MetricBase.MetricValues)

	backend.GetRawData(name, from, to, out)
	for m := range out {
		list = append(list, m)
	}

	return list
}
