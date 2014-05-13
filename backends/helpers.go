package backends

import (
	"github.com/mb0/glob"
	"github.com/msiebuhr/MetricBase/metrics"
)

// GetMetricsAsList fetches all the metrics in a backend and returns it as a
// list.
func GetMetricsAsList(backend Backend) []string {
	list := make([]string, 0)
	out := make(chan string)

	backend.GetMetricsList(out)
	for name := range out {
		list = append(list, name)
	}
	return list
}

// Apply glob pattern to metrics from a backend
func GlobMetricsAsList(pattern string, backend Backend) ([]string, error) {
	conf := glob.Default()
	conf.Separator = '.'
	globber, err := glob.New(conf)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0)
	out := make(chan string)

	backend.GetMetricsList(out)
	for name := range out {
		match, err := globber.Match(pattern, name)
		if err != nil {
			return nil, err
		}

		if match {
			list = append(list, name)
		}
	}
	return list, nil
}

// GetDataAsList fetches the relevant data and returns it as a list.
func GetDataAsList(backend Backend, name string, from, to int64) []metrics.MetricValue {
	list := make([]metrics.MetricValue, 0)
	out := make(chan metrics.MetricValue)

	backend.GetRawData(name, from, to, out)
	for m := range out {
		list = append(list, m)
	}

	return list
}
