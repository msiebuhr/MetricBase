package backends

import (
	"time"

	"github.com/msiebuhr/MetricBase/metrics"
)

type Backend interface {
	// Start the backend. Often in a go-routine
	Start()
	// Exit the Start()'ed go-routine.
	Stop()

	// Add the given metric to the backend
	AddMetricChan(chan metrics.Metric)
	// Return a stream of data for the givem name, from and to arguments
	GetRawData(string, time.Time, time.Time, chan metrics.MetricValue)
	// Get a list of the metrics the backend knows about
	GetMetricsList(chan string)
}
