package backends

import (
	"github.com/msiebuhr/MetricBase"
)

// Commonly used internal data structures
type dataRequest struct {
	Name   string
	From   int64
	To     int64
	Result chan MetricBase.MetricValues
}
