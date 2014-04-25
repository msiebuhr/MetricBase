package query

import (
	"github.com/msiebuhr/MetricBase"
)

/* Request data from a built query tree
 *
 * TODO: Should it include some sort of desired resolution? That
 * could make things like culling data up-front and various
 * aggregators more efficient.
 */
type Request struct {
	Backend  MetricBase.Backend
	From, To int64
}

func NewRequest(from, to int64) Request {
	return Request{
		From: from,
		To:   to,
	}
}
