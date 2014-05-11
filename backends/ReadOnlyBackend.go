package backends

import (
	"github.com/msiebuhr/MetricBase/metrics"
)

type ReadOnlyBackend struct {
	data map[string][]metrics.MetricValue

	stopChan     chan bool
	listRequests chan chan string
	dataRequests chan dataRequest
}

func NewReadOnlyBackend(data ...*metrics.Metric) *ReadOnlyBackend {
	r := &ReadOnlyBackend{
		data:         make(map[string][]metrics.MetricValue),
		stopChan:     make(chan bool),
		listRequests: make(chan chan string),
		dataRequests: make(chan dataRequest),
	}

	// Add given metrics to backend
	for _, m := range data {
		r.data[m.Name] = append(
			r.data[m.Name],
			metrics.MetricValue{Time: m.Time, Value: m.Value},
		)
	}

	return r
}

func (r *ReadOnlyBackend) Start() {
	go func() {
		for {
			select {
			case req := <-r.listRequests:
				for key := range r.data {
					req <- key
				}
				close(req)
			case req := <-r.dataRequests:
				if _, ok := r.data[req.Name]; ok {
					for _, data := range r.data[req.Name] {
						if req.From < data.Time && data.Time < req.To {
							req.Result <- data
						}
					}
				}
				close(req.Result)
			case <-r.stopChan:
				close(r.stopChan)
				close(r.listRequests)
				close(r.dataRequests)
				return
			}
		}
	}()
}

// AddMetric is a dummy function that throws all given data away.
func (r *ReadOnlyBackend) AddMetric(c metrics.Metric) {}

func (r *ReadOnlyBackend) Stop() { r.stopChan <- true }

func (r *ReadOnlyBackend) GetMetricsList(results chan string) {
	r.listRequests <- results
}

func (r *ReadOnlyBackend) GetRawData(name string, from, to int64, result chan metrics.MetricValue) {
	r.dataRequests <- dataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}
