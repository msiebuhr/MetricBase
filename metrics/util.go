package metrics

import (
	"time"
)

// TimeToUint40 encodes a Unix-timestamp with second-resolution and
// encodes it in five bytes.
//
// Behavior is undefined for dates before 1970 and somewhere after year 2300.
func TimeToUint40(b []byte, t time.Time) {
	v := uint64(t.Unix())
	b[0] = byte(v >> 32)
	b[1] = byte(v >> 24)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 8)
	b[4] = byte(v)
}

func Uint40ToTime(b []byte) time.Time {
	i := int64(b[4]) |
		int64(b[3])<<8 |
		int64(b[2])<<16 |
		int64(b[1])<<24 |
		int64(b[0])<<32
	return time.Unix(i, 0)
}
