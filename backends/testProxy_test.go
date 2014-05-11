package backends

import (
	"testing"

	"github.com/msiebuhr/MetricBase/metrics"
)

func TestTestProxy(t *testing.T) {
	backend := NewTestProxy(NewReadOnlyBackend(
		metrics.NewMetric("foo", 1, 1),
		metrics.NewMetric("foo", 2, 2),
	))

	// Start backend
	backend.Start()
	defer backend.Stop()

	// Read back the data
	data := GetDataAsList(backend, "foo", 0, 5)
	if len(data) != 2 {
		t.Fatalf("Expected to get two results, got %d", len(data))
	}

	if data[0].Value != 1 {
		t.Errorf("Expected data[0].Value=1, got '%f'.", data[0].Value)
	}
	if data[0].Time != 1 {
		t.Errorf("Expected data[0].Time=1, got '%d'.", data[0].Time)
	}

	if data[1].Value != 2 {
		t.Errorf("Expected data[0].Value=2, got '%f'.", data[1].Value)
	}
	if data[1].Time != 2 {
		t.Errorf("Expected data[0].Time=2, got '%d'.", data[1].Time)
	}

	data = GetDataAsList(backend, "foo", -1, 0)
	if len(data) != 0 {
		t.Errorf("Expected no data, got %v", data)
	}

	// Get some test.X data
	data = GetDataAsList(backend, "test.sin", 1, 100)
	if len(data) == 0 {
		t.Errorf("Expected some data from test.sin, got none")
	}
	//t.Errorf("%#v", data)
}
