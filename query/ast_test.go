package query

import (
	"testing"

	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/metrics"
)

func TestStringer(t *testing.T) {
	var stringerTests = []struct {
		query    string
		stringed string
	}{
		// Negative tests
		//{"", "[]"},
		//{"foo", "[]"},

		// Matches
		{"foo.bar", "foo.bar"},
		{"foo.*", "foo.*"},

		// Functions
		{"scale(foo.*, 10.0)", "scale(foo.*, 10)"},
		{"scale(42,foo.*)", "scale(foo.*, 42)"},
		{"scale(foo.*)", "scale(foo.*, 1)"},
		{"scale(scale(foo.*, 10), 5)", "scale(scale(foo.*, 10), 5)"},

		// Should err
		//{`scale("should")`, ""},
		//{"scale(foo.*, 10, 5)", "scale(foo.*, 15)"},
	}

	for _, tt := range stringerTests {
		// Build a query
		out, err := ParseGraphiteQuery(tt.query)

		if err != nil {
			t.Errorf("Err'd on input %s: %v", tt.query, err)
			continue
		}

		stringed := out.String()
		if stringed != tt.stringed {
			t.Errorf("Expected String on '%v' to return '%v', got '%v'.", tt.query, tt.stringed, stringed)
		}
	}
}

func TestSimpleQuery(t *testing.T) {
	backend := backends.NewReadOnlyBackend(
		metrics.NewMetric("foo.bar", 3.14, 10),
		metrics.NewMetric("foo.bar", 42.0, 20),
		metrics.NewMetric("foo.baz", 10, 10),
		metrics.NewMetric("foo.baz", 20, 20),
	)
	backend.Start()
	defer backend.Stop()

	var stringerTests = []struct {
		query    string
		from, to int64
		output   []Response
	}{
		// Should match everything
		{"foo.bar", 0, 30, []Response{}},

		// Should match nothing
		//{"foo.bar", 0, 10, []Response{}},

		// Just a number
		{"42", 0, 1, []Response{}},

		// Funcntions
		{"scale(foo.bar, 10.1)", 0, 30, []Response{}},
		{"scale(foo.bar, foo.baz, 10)", 0, 30, []Response{}},
	}

	for _, tt := range stringerTests {
		// Build a query
		out, err := ParseGraphiteQuery(tt.query)

		if err != nil {
			t.Errorf("Err'd on input %s: %v", tt.query, err)
			continue
		}

		res, err := out.Query(Request{backend, tt.from, tt.to})
		if err != nil {
			t.Errorf("Query() errored: %v", err)
			continue
		}

		// Expect to get some data back
		for i := range res {
			data := res[i].GetAllMetrics()
			if len(data) == 0 {
				t.Errorf("Expected query '%s' [%d, %d] to return some data.", out, tt.from, tt.to)
			}
		}
	}
}

func TestParseError(t *testing.T) {
	var errorTests = []struct {
		query string
		err   string
	}{
		{`""`, "error at position 1"},
	}

	for _, tt := range errorTests {
		// Build a query
		_, err := ParseGraphiteQuery(tt.query)

		if err == nil {
			t.Errorf("Expected error on input '%v'.", tt.query)
			continue
		}

		if err.Error() != tt.err {
			t.Errorf("Expected query '%v' to return error '%v', got '%v'", tt.query, tt.err, err)
		}
	}
}
