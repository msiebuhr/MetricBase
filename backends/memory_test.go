package backends

import (
	"testing"
	"time"

	"github.com/msiebuhr/MetricBase/metrics"
)

func generateTestStoreAndGet(backend Backend, t *testing.T) {
	// Start backend
	backend.Start()
	defer backend.Stop()

	// Load some data and read it back out
	backend.AddMetric(*metrics.NewMetric("foo.bar", 3.14, 100))

	time.Sleep(time.Millisecond)

	// Read back list of metrics
	metricNames := GetMetricsAsList(backend)
	if len(metricNames) != 1 {
		t.Errorf("Expected to get one metric back, got %d", len(metricNames))
	} else if metricNames[0] != "foo.bar" {
		t.Errorf("Expected the metric name to be 'foo.bar', got '%v'", metricNames[0])
	}

	// Read back the data
	data := GetDataAsList(backend, "foo.bar", 99, 101)
	if len(data) != 1 {
		t.Fatalf("Expected to get one result, got %d", len(data))
	}

	if data[0].Value != 3.14 {
		t.Errorf("Expected data[0].Value=3.14, got '%f'.", data[0].Value)
	}
	if data[0].Time != 100 {
		t.Errorf("Expected data[0].Time=100, got '%d'.", data[0].Time)
	}

	// Make a query that shouldn't return any data
	data = GetDataAsList(backend, "foo.bar", 0, 0)
	if len(data) != 0 {
		t.Errorf("Expected no data for query (foo.bar, 0, 0), got %v", data)
	}
}

func TestMemoryStoreAndGet(t *testing.T) {
	// Create tempdir (& remove afterwards)
	generateTestStoreAndGet(NewMemoryBackend(), t)
}
