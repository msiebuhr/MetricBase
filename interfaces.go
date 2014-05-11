package MetricBase

import (
	"github.com/msiebuhr/MetricBase/metrics"
)

type Backend interface {
	Start()
	Stop()
	AddMetric(metrics.Metric)
	GetRawData(string, int64, int64, chan metrics.MetricValue)
	GetMetricsList(chan string)
}

type Frontend interface {
	SetBackend(Backend)
	Start()
	Stop()
}
