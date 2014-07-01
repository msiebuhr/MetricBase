package testProxy

import (
	"testing"
	"time"

	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/backends/readOnly"
	"github.com/msiebuhr/MetricBase/metrics"
)

func TestTestProxy(t *testing.T) {
	backend := NewTestProxy(readOnly.NewReadOnlyBackend(
		metrics.NewMetric("foo", 1, 1),
		metrics.NewMetric("foo", 2, 2),
	))

	// Start backend
	backend.Start()
	defer backend.Stop()

	// Read back the data
	data := backends.GetDataAsList(backend, "foo", time.Unix(0, 0), time.Unix(5, 0))
	if len(data) != 2 {
		t.Fatalf("Expected to get two results, got %d", len(data))
	}

	if data[0].Value != 1 {
		t.Errorf("Expected data[0].Value=1, got '%f'.", data[0].Value)
	}
	if data[0].Time != time.Unix(1, 0) {
		t.Errorf("Expected data[0].Time=1, got '%d'.", data[0].Time)
	}

	if data[1].Value != 2 {
		t.Errorf("Expected data[0].Value=2, got '%f'.", data[1].Value)
	}
	if data[1].Time != time.Unix(2, 0) {
		t.Errorf("Expected data[0].Time=2, got '%d'.", data[1].Time)
	}

	data = backends.GetDataAsList(backend, "foo", time.Unix(-1, 0), time.Unix(0, 0))
	if len(data) != 0 {
		t.Errorf("Expected no data, got %v", data)
	}

	// Get some test.X data
	data = backends.GetDataAsList(backend, "test.sin", time.Unix(1, 0), time.Unix(100, 0))
	if len(data) == 0 {
		t.Errorf("Expected some data from test.sin, got none")
	}
	//t.Errorf("%#v", data)
}
