package backends

import (
	"github.com/msiebuhr/MetricBase"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func generateTestStoreAndGet(backend MetricBase.Backend, t *testing.T) {
	// Start backend
	backend.Start()
	defer backend.Stop()

	// Load some data and read it back out
	addChan := make(chan MetricBase.Metric, 10)
	backend.AddMetrics(addChan)

	m := MetricBase.Metric{}
	m.Name = "foo.bar"
	m.Time = 100
	m.Value = 3.14
	addChan <- m
	close(addChan)

	time.Sleep(time.Millisecond)

	// Read back list of metrics
	metricNameChan := make(chan string, 1)
	backend.GetMetricsList(metricNameChan)
	for metric := range metricNameChan {
		if metric != "foo.bar" {
			t.Errorf("Expected to get metric 'foo.bar', got '%v'", metric)
		}
	}

	// Read back the data
	metricChan := make(chan MetricBase.MetricValues)
	backend.GetRawData("foo.bar", 0, 0, metricChan)
	for data := range metricChan {
		if data.Time != 100 {
			t.Errorf("Expected to get time '100', got '%v'", data.Time)
		}
		if data.Value != 3.14 {
			t.Errorf("Expected to get time '3.14', got '%v'", data.Value)
		}
	}
}

func TestLevelDbStoreAndGet(t *testing.T) {
	// Create tempdir (& remove afterwards)
	dir, err := ioutil.TempDir(".", "tmp_storage_test")
	if err != nil {
		t.Fatalf("Could not create tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	generateTestStoreAndGet(NewLevelDb(dir), t)
}

func TestMemoryStoreAndGet(t *testing.T) {
	// Create tempdir (& remove afterwards)
	generateTestStoreAndGet(NewMemoryBackend(), t)
}
