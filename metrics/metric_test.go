package metrics

import "testing"

// String()
func TestMetricString(t *testing.T) {
	m := NewMetric("foo", 10.1, 42)
	s := m.String()

	if s != "foo 10.1 42" {
		t.Errorf("Expected 'foo 10.1 42', got '%v'", s)
	}
}

func BenchmarkMetricString(b *testing.B) {
	m := NewMetric("foo.bar", 1.23, 123)
	for i := 0; i < b.N; i++ {
		s := m.String()
		b.SetBytes(int64(len(s)))
	}
}
