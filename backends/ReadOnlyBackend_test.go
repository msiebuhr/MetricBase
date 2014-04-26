package backends

import (
	"testing"

	"github.com/msiebuhr/MetricBase"
)

func TestReadOnlyBackend(t *testing.T) {
	backend := NewReadOnlyBackend(
		MetricBase.NewMetric("foo", 1, 1),
		MetricBase.NewMetric("foo", 2, 2),
	)

	// Start backend
	backend.Start()
	defer backend.Stop()

	// Read back list of metrics
	metricNames := GetMetricsAsList(backend)
	if len(metricNames) != 1 {
		t.Errorf("Expected to get one metric back, got %d", len(metricNames))
	} else if metricNames[0] != "foo" {
		t.Errorf("Expected the metric name to be 'foo', got '%v'", metricNames[0])
	}

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

}
