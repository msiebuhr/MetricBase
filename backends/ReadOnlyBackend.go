package backends

import (
	"github.com/msiebuhr/MetricBase"
)

type ReadOnlyBackend struct {
	data map[string][]MetricBase.MetricValues

	stopChan     chan bool
	listRequests chan chan string
	dataRequests chan dataRequest
}

func NewReadOnlyBackend(metrics ...*MetricBase.Metric) *ReadOnlyBackend {
	r := &ReadOnlyBackend{
		data:         make(map[string][]MetricBase.MetricValues),
		stopChan:     make(chan bool),
		listRequests: make(chan chan string),
		dataRequests: make(chan dataRequest),
	}

	// Add given metrics to backend
	for _, m := range metrics {
		r.data[m.Name] = append(
			r.data[m.Name],
			MetricBase.MetricValues{Time: m.Time, Value: m.Value},
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
func (r *ReadOnlyBackend) AddMetric(c MetricBase.Metric) {}

func (r *ReadOnlyBackend) Stop() { r.stopChan <- true }

func (r *ReadOnlyBackend) GetMetricsList(results chan string) {
	r.listRequests <- results
}

func (r *ReadOnlyBackend) GetRawData(name string, from, to int64, result chan MetricBase.MetricValues) {
	r.dataRequests <- dataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}
