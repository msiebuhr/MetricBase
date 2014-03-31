package frontends

import (
	"testing"
	//"time"
	"github.com/msiebuhr/MetricBase"
)

// What about some basic lines
func TestSingleLineParsing(t *testing.T) {
	t.Parallel()
	var linetests = []struct {
		in  string
		out MetricBase.Metric
	}{
		{
			"foo 1 2",
			*MetricBase.NewMetric("foo", 1, 2),
		},
		{
			"a.b.c 4.2 42",
			*MetricBase.NewMetric("a.b.c", 4.2, 42),
		},
	}

	for i, tt := range linetests {
		_, outMetric := parseGraphiteLine(tt.in)
		if outMetric != tt.out {
			t.Errorf("%d. parseGraphiteLine(%s) => %s, want %s", i, tt.in, outMetric, tt.out)
		}
	}
}

func TestSingleLineParserFail(t *testing.T) {
	t.Parallel()
	failLines := []string{
		"",
		"one",
		"two elements",
		"four elements to go",
		"foo 1 should_be_int",
		"name should_be_float 1",
	}

	for i, tt := range failLines {
		err, _ := parseGraphiteLine(tt)
		if err == nil {
			t.Errorf("%d. Expected parseGraphiteLine(%s) to return error", i, tt)
		}
	}
}
