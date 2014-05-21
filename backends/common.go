package backends

import (
	"time"

	"github.com/msiebuhr/MetricBase/metrics"
)

// Commonly used internal data structures
type dataRequest struct {
	Name   string
	From   time.Time
	To     time.Time
	Result chan metrics.MetricValue
}
