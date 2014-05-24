package metrics

import (
	"testing"
	"time"
)

func TestTimeToUint40(t *testing.T) {
	var tests = []struct {
		in  time.Time
		out []byte
	}{
		// What about zero
		{time.Unix(0, 0), []byte{0, 0, 0, 0, 0}},

		// Discards milliseconds
		{time.Unix(1, 1), []byte{0, 0, 0, 0, 1}},

		// Larger numbers
		{time.Unix(255, 1), []byte{0, 0, 0, 0, 0xff}},
		{time.Unix(1<<8, 1), []byte{0, 0, 0, 1, 0}},
		{time.Unix(1<<16, 1), []byte{0, 0, 1, 0, 0}},
		{time.Unix(1<<24, 1), []byte{0, 1, 0, 0, 0}},
		{time.Unix(1<<32, 1), []byte{1, 0, 0, 0, 0}},

		// Shoot over
		{time.Unix(1<<40, 1), []byte{0, 0, 0, 0, 0}},
	}

	for _, tt := range tests {
		outArray := make([]byte, 5)
		TimeToUint40(outArray, tt.in)

		if len(outArray) != len(tt.out) ||
			outArray[0] != tt.out[0] ||
			outArray[1] != tt.out[1] ||
			outArray[2] != tt.out[2] ||
			outArray[3] != tt.out[3] ||
			outArray[4] != tt.out[4] {
			t.Errorf("Expected %v to encode to %v, got %v.", tt.in, tt.out, outArray)
		}
	}
}

// Should probably use testing/quick, but I can't quite get it to run.
func TestTimeAndUint40BackAndForth(t *testing.T) {
	buf := make([]byte, 5)
	var i int64
	for i = 0; i < 1e10; i += 1e4 {
		in := time.Unix(i, 0)
		TimeToUint40(buf, in)
		out := Uint40ToTime(buf)

		if !in.Equal(out) {
			t.Errorf("Encoding %v -> %v -> %v fails", in, buf, out)
		}
	}
}

func BenchmarkTimeToUint40(b *testing.B) {
	out := make([]byte, 5)
	t := time.Unix(int64(b.N), 0)
	for i := 0; i < b.N; i++ {
		TimeToUint40(out, t)
	}
}

func TestUint40ToTime(t *testing.T) {
	var tests = []struct {
		in  []byte
		out time.Time
	}{
		{[]byte{0, 0, 0, 0, 0}, time.Unix(0, 0)},

		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, time.Date(36812, 2, 20, 0, 36, 15, 0, time.UTC)},
	}

	for _, tt := range tests {
		res := Uint40ToTime(tt.in)

		if !res.Equal(tt.out) {
			t.Errorf("Expected %v to decode to %v, got %v.", tt.in, tt.out, res)
		}
	}
}

func BenchmarkUint40ToTime(b *testing.B) {
	in := make([]byte, 5)
	for i := 0; i < b.N; i++ {
		_ = Uint40ToTime(in)
	}
}
