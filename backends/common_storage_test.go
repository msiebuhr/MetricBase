package backends

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/msiebuhr/MetricBase/backends/boltdb"
	"github.com/msiebuhr/MetricBase/backends/memory"
	"github.com/msiebuhr/MetricBase/metrics"
)

func generateTestStoreAndGet(backend Backend, t *testing.T) {
	// Start backend
	backend.Start()
	defer backend.Stop()

	// Load some data and read it back out
	inChan := make(chan metrics.Metric)
	defer close(inChan)
	backend.AddMetricChan(inChan)
	inChan <- *metrics.NewMetric("foo.bar", 3.14, 100)

	time.Sleep(time.Millisecond)

	// Read back list of metrics
	metricNames := GetMetricsAsList(backend)
	if len(metricNames) != 1 {
		t.Errorf("Expected to get one metric back, got %d", len(metricNames))
	} else if metricNames[0] != "foo.bar" {
		t.Errorf("Expected the metric name to be 'foo.bar', got '%v'", metricNames[0])
	}

	// Read back the data
	data := GetDataAsList(backend, "foo.bar", time.Unix(99, 0), time.Unix(101, 0))
	if len(data) != 1 {
		t.Fatalf("Expected to get one result, got %d", len(data))
	}

	if data[0].Value != 3.14 {
		t.Errorf("Expected data[0].Value=3.14, got '%f'.", data[0].Value)
	}
	if data[0].Time != time.Unix(100, 0) {
		t.Errorf("Expected data[0].Time=100, got '%d'.", data[0].Time)
	}

	// Make a query that shouldn't return any data
	data = GetDataAsList(backend, "foo.bar", time.Unix(0, 0), time.Unix(0, 0))
	if len(data) != 0 {
		t.Errorf("Expected no data for query (foo.bar, 0, 0), got %v", data)
	}
}

func TestMemoryStoreAndGet(t *testing.T) {
	// Create tempdir (& remove afterwards)
	generateTestStoreAndGet(memory.NewMemoryBackend(), t)
}

func TestBoltStoreAndGet(t *testing.T) {
	f, err := ioutil.TempFile(".", "boltdb_test")
	if err != nil {
		t.Fatalf("Could not create test directory.")
	}
	defer os.Remove(f.Name())

	db, err := boltdb.NewBoltBackend(f.Name())
	if err != nil {
		t.Fatalf("Could not start BoltDB %v", err)
	}
	generateTestStoreAndGet(db, t)
}
