package query

import (
	"time"

	"github.com/msiebuhr/MetricBase/backends"
)

/* Request data from a built query tree
 *
 * TODO: Should it include some sort of desired resolution? That
 * could make things like culling data up-front and various
 * aggregators more efficient.
 */
type Request struct {
	Backend  backends.Backend
	From, To time.Time
}

func NewRequest(from, to time.Time) Request {
	return Request{
		From: from,
		To:   to,
	}
}
