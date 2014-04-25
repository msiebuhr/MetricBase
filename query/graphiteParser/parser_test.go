package graphiteParser

import (
	"testing"
)

func TestSimpleStrigs(t *testing.T) {
	var tests = []struct {
		query    string
		stringed string
	}{
		// Negative tests
		//{"", "[]"},
		//{"foo", "[]"},

		// Matches
		{"foo.bar", "foo.bar"},
		{"foo.*", "foo.*"},

		// Node types
		{"3.14", "3.14"},
		{"42", "42"},
		{"\"string\"", "\"string\""},

		// Functions
		{"runningAverage( foo.* , 10 )", "runningAverage(foo.*, 10)"},
	}

	for _, tt := range tests {
		// Build a query
		out, err := Parse(tt.query)

		if err != nil {
			t.Errorf("Didn't expect an error for input %v.", tt.query)
			continue
		}

		stringed := out.String()
		if stringed != tt.stringed {
			t.Errorf("Expected String on '%v' to return '%v', got '%v'.", tt.query, tt.stringed, stringed)
		}
	}
}
