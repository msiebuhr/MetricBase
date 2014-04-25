package query

import (
	"github.com/msiebuhr/MetricBase"
)

type Response struct {
	Meta map[string]string
	Data chan MetricBase.MetricValues
}

func (r Response) GetAllMetrics() []MetricBase.MetricValues {
	out := make([]MetricBase.MetricValues, 0)
	for m := range r.Data {
		out = append(out, m)
	}
	return out
}

func NewResponse() Response {
	return Response{
		Meta: make(map[string]string),
		Data: make(chan MetricBase.MetricValues, 100),
	}
}
