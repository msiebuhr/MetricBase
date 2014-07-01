/* Loads internally exposed data
 *
 * I would love to go through expvar, but there doesn't seem to be a sane way
 * without somehow parsing the strings and what not.
 */
package internalMetrics

import (
	"runtime"
	"time"

	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/metrics"
)

type internalMetrics struct {
	backend  backends.Backend
	stopChan chan bool
	interval time.Duration
}

func NewInternalMetrics(interval time.Duration) *internalMetrics {
	// TODO: Make sure interval it at least one second
	return &internalMetrics{
		interval: interval,
		stopChan: make(chan bool),
		backend:  nil,
	}
}

func (e *internalMetrics) SetBackend(backend backends.Backend) {
	e.backend = backend
}

// Submit metrics to the given channel
func submitUint64Metric(name string, when time.Time, value uint64) metrics.Metric {
	return metrics.Metric{
		Name: name,
		MetricValue: metrics.MetricValue{
			Value: float64(value),
			Time:  when,
		},
	}
}

func (e *internalMetrics) Start() {
	go func() {
		// Ticker
		t := time.NewTicker(e.interval)

		// Output channel
		outChan := make(chan metrics.Metric, 100)
		e.backend.AddMetricChan(outChan)

		// Misc structures
		m := &runtime.MemStats{}

		for {
			select {
			case when := <-t.C: // Receive tick
				when = when.Round(time.Second).UTC()
				// Fetch runtime.Memstats
				runtime.ReadMemStats(m)

				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.Alloc", when, m.Alloc)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.TotalAlloc", when, m.TotalAlloc)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.Lookups", when, m.Lookups)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.Mallocs", when, m.Mallocs)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.Frees", when, m.Frees)

				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.HeapAlloc", when, m.HeapAlloc)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.HeapSys", when, m.HeapSys)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.HeapIdle", when, m.HeapIdle)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.HeapInuse", when, m.HeapInuse)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.HeapReleased", when, m.HeapReleased)
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.HeapObjects", when, m.HeapObjects)

				// Do some intelligent stats on GC stuff
				outChan <- submitUint64Metric("MetricBase.runtime.MemStats.GC.LastPause", when, m.PauseNs[(m.NumGC+255)%256])

				// Misc stuff
				outChan <- submitUint64Metric("MetricBase.runtime.NumGoroutine", when, uint64(runtime.NumGoroutine()))
				outChan <- submitUint64Metric("MetricBase.runtime.GOMAXPROCS", when, uint64(runtime.GOMAXPROCS(0)))

			case <-e.stopChan: // Stop
				t.Stop()
				close(outChan)
				close(e.stopChan)
				return
			}
		}
	}()
}

func (e *internalMetrics) Stop() {
	e.stopChan <- true
}
