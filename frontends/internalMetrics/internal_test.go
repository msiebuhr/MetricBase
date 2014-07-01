package internalMetrics

import (
	"testing"
	"time"

	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/backends/memory"
)

func TestExpvarLoader(t *testing.T) {
	f := NewInternalMetrics(time.Millisecond)
	b := memory.NewMemoryBackend()
	f.SetBackend(b)

	b.Start()
	f.Start()

	time.Sleep(time.Millisecond * 2)

	f.Stop()

	// Verify that b got some metrics
	names := backends.GetMetricsAsList(b)
	if len(names) == 0 {
		t.Fatalf("Expected some metrics back, got none")
	}

	t.Logf("Got metrics %v", names)

	// Read some of the data
	for _, name := range names {
		values := backends.GetDataAsList(b, name, time.Now().Add(-1*time.Second), time.Now().Add(time.Second))

		t.Logf("%s: %v", name, values)
	}
}
