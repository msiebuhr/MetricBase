package backends

import (
	"math"
	"strings"

	"github.com/msiebuhr/MetricBase"
)

var patterns map[string]func(i int64) float64

func init() {
	patterns = make(map[string]func(i int64) float64)

	// Basic math stuff
	patterns["sin.hour"] = func(i int64) float64 {
		delta := 60.0 * 60
		return math.Sin(2 * math.Pi / delta * float64(i))
	}
	patterns["sin.day"] = func(i int64) float64 {
		delta := 60.0 * 60 * 24
		return math.Sin(2 * math.Pi / delta * float64(i))
	}
	patterns["sin.week"] = func(i int64) float64 {
		delta := 60.0 * 60 * 24 * 7
		return math.Sin(2 * math.Pi / delta * float64(i))
	}

	// Constants
	patterns["const.1"] = func(i int64) float64 { return 1.0 }

}

// TestProxy sits in front of another backend and services various kinds of
// dummy data under the 'test.'-prefix.
type TestProxy struct {
	nextBackend MetricBase.Backend
}

func NewTestProxy(next MetricBase.Backend) *TestProxy {
	return &TestProxy{nextBackend: next}
}

func (tp *TestProxy) Start() { tp.nextBackend.Start() }
func (tp *TestProxy) Stop()  { tp.nextBackend.Stop() }

func (tp *TestProxy) AddMetric(metric MetricBase.Metric) {
	tp.nextBackend.AddMetric(metric)
}

func (tp *TestProxy) GetMetricsList(results chan string) {
	go func() {
		// Inject known test metrics
		for name := range patterns {
			results <- "test." + name
		}

		// Ask whatever backend we're running on about the rest
		// Don't close results, as callee will do that:w
		tp.nextBackend.GetMetricsList(results)
	}()
}

func (tp *TestProxy) GetRawData(name string, from, to int64, result chan MetricBase.MetricValues) {
	if !strings.HasPrefix(name, "test.") {
		tp.nextBackend.GetRawData(name, from, to, result)
		return
	}

	// Default function to use for generating data
	f := func(i int64) float64 { return float64(i) }

	if fun, ok := patterns[name[5:]]; ok {
		f = fun
	}

	// Run the generator-function on relevant input
	go func() {
		defer close(result)
		// Generate about 4000 points, so we don't generate too little or too much data, pending the resolution
		delta := (to - from) / 4000
		if delta <= 0 {
			delta = 1
		}
		for i := from; i <= to; i += delta {
			result <- MetricBase.MetricValues{
				Time:  i,
				Value: f(i),
			}
		}
	}()
}