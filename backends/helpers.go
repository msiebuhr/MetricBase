package backends

import (
	"github.com/msiebuhr/MetricBase"
	"github.com/msiebuhr/MetricBase/metrics"
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
func GetDataAsList(backend MetricBase.Backend, name string, from, to int64) []metrics.MetricValue {
	list := make([]metrics.MetricValue, 0)
	out := make(chan metrics.MetricValue)

	backend.GetRawData(name, from, to, out)
	for m := range out {
		list = append(list, m)
	}

	return list
}
