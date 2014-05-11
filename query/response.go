package query

import (
	"github.com/msiebuhr/MetricBase/metrics"
)

type Response struct {
	Meta map[string]string
	Data chan metrics.MetricValue
}

func (r Response) GetAllMetrics() []metrics.MetricValue {
	out := make([]metrics.MetricValue, 0)
	for m := range r.Data {
		out = append(out, m)
	}
	return out
}

func NewResponse() Response {
	return Response{
		Meta: make(map[string]string),
		Data: make(chan metrics.MetricValue, 100),
	}
}
