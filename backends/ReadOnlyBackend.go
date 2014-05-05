package backends

import (
	"github.com/msiebuhr/MetricBase"
)

type ReadOnlyBackend struct {
	data map[string][]MetricBase.MetricValues

	stopChan        chan bool
	listRequestChan chan chan string
	dataRequestChan chan dataRequest
}

func NewReadOnlyBackend(metrics ...*MetricBase.Metric) *ReadOnlyBackend {
	r := &ReadOnlyBackend{
		data:            make(map[string][]MetricBase.MetricValues),
		stopChan:        make(chan bool),
		listRequestChan: make(chan chan string),
		dataRequestChan: make(chan dataRequest),
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
			case req := <-r.listRequestChan:
				for key := range r.data {
					req <- key
				}
				close(req)
			case req := <-r.dataRequestChan:
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
				close(r.listRequestChan)
				close(r.dataRequestChan)
				return
			}
		}
	}()
}

// TODO: Dummy addmetrics function
func (r *ReadOnlyBackend) AddMetrics(c chan MetricBase.Metric) {
	// Throw all data away
	go func() {
		for _ = range c {
		}
	}()
}

func (r *ReadOnlyBackend) Stop() { r.stopChan <- true }

func (r *ReadOnlyBackend) GetMetricsList(results chan string) {
	r.listRequestChan <- results
}

func (r *ReadOnlyBackend) GetRawData(name string, from, to int64, result chan MetricBase.MetricValues) {
	r.dataRequestChan <- dataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}
