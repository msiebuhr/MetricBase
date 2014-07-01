package backends

import (
	"testing"

	"github.com/msiebuhr/MetricBase/backends/readOnly"
	"github.com/msiebuhr/MetricBase/metrics"
)

func TestGlobHelper(t *testing.T) {
	backend := readOnly.NewReadOnlyBackend(
		metrics.NewMetric("foo", 1, 1),
		metrics.NewMetric("foo.1.bar", 2, 2),
		metrics.NewMetric("foo.a2.bar", 2, 2),
	)

	// Start backend
	backend.Start()
	defer backend.Stop()

	var globberTests = []struct {
		glob string
		out  []string
	}{
		// Nothing
		{"foo.*", []string{}},
		{"foo.bar", []string{}},

		// Globless
		{"foo", []string{"foo"}},
		{"foo.1.bar", []string{"foo.1.bar"}},

		// Sub matches
		{"foo.*.bar", []string{"foo.1.bar", "foo.a2.bar"}},
		{"foo.?.bar", []string{"foo.1.bar"}},
		{"foo.a*.bar", []string{"foo.a2.bar"}},
		{"foo.a?.bar", []string{"foo.a2.bar"}},
	}

	for _, tt := range globberTests {
		// Build a query
		out, err := GlobMetricsAsList(tt.glob, backend)

		if err != nil {
			t.Errorf("Err'd on input %s: %v", tt.glob, err)
			continue
		}

		if len(out) != len(tt.out) {
			t.Errorf("Expected glob %v to return %v, got %v.", tt.glob, tt.out, out)
		}

		// Put all requred names in a map and remove all available elements
		m := make(map[string]bool)
		for _, s := range tt.out {
			m[s] = false
		}

		for _, s := range out {
			if _, ok := m[s]; !ok {
				t.Errorf("Expected glob %v to include %v. It didn't.", tt.glob, s)
			}

			if m[s] {
				t.Errorf("Glob %v returned duplicate element %v.", tt.glob, s)
			}

			m[s] = true
		}
	}
}

func TestGlobPatternPrefix(t *testing.T) {
	var globberTests = []struct {
		glob   string
		prefix string
	}{
		{"foo.bar", "foo.bar"},

		{"foo.*.*.bar", "foo."},
		{"foo?", "foo"},
		{"foo]bar", "foo"},
		{"foo[bar", "foo"},

		{"trailing\\", "trailing\\"},
		{"escaped.\\*.star", "escaped.\\*.star"},
	}

	for _, tt := range globberTests {
		// Build a query
		out := GlobPatternPrefix(tt.glob)

		if out != tt.prefix {
			t.Errorf("Expected '%v' to have prefix '%v', got '%v'.", tt.glob, tt.prefix, out)
		}
	}
}
