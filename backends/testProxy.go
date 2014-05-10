package backends

import (
	"math"
	"strings"

	"github.com/msiebuhr/MetricBase"
)

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
		//defer close(results)
		// Inject known test metrics
		results <- "test.sin"

		// Ask whatever backend we're running on about the rest
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

	// Switch on what the rest of the name is
	switch name[5:] {
	case "sin":
		f = func(i int64) float64 {
			return math.Sin(float64(i))
		}
	}

	// Run the generator-function on relevant input
	go func() {
		defer close(result)
		// Generate sine curve
		for i := from; i <= to; i += 10 {
			result <- MetricBase.MetricValues{
				Time:  i,
				Value: f(i),
			}
		}
	}()
}
