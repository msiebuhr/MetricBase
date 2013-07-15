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
	addReq := &MetricBase.AddRequest{Data: make(chan MetricBase.Metric, 100)}
	backend.Add(*addReq)

	m := MetricBase.Metric{}
	m.Name = "foo.bar"
	m.Time = 100
	m.Value = 3.14
	addReq.Data <- m
	close(addReq.Data)

	time.Sleep(time.Millisecond)

	// Read back list of metrics
	listReq := &MetricBase.ListRequest{Result: make(chan string, 1)}
	backend.List(*listReq)
	for metric := range listReq.Result {
		if metric != "foo.bar" {
			t.Errorf("Expected to get metric 'foo.bar', got '%v'", metric)
		}
	}

	// Read back the data
	dataReq := &MetricBase.DataRequest{
		Name:   "foo.bar",
		Result: make(chan MetricBase.MetricValues),
	}

	// Fetch the data
	backend.Data(*dataReq)
	for data := range dataReq.Result {
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

	generateTestStoreAndGet(CreateLevelDb(dir), t)
}

func TestMemoryStoreAndGet(t *testing.T) {
	// Create tempdir (& remove afterwards)
	generateTestStoreAndGet(CreateMemoryBackend(), t)
}
