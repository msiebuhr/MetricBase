package backends

import (
	"github.com/msiebuhr/MetricBase"
)

type MemoryBackend struct {
	data map[string][]MetricBase.MetricValues

	stopChan     chan bool
	addChan      chan MetricBase.Metric
	listRequests chan chan string
	dataRequests chan dataRequest
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		data:         make(map[string][]MetricBase.MetricValues),
		stopChan:     make(chan bool),
		addChan:      make(chan MetricBase.Metric),
		listRequests: make(chan chan string),
		dataRequests: make(chan dataRequest),
	}
}

func (m *MemoryBackend) Start() {
	go func() {
		for {
			select {
			case metric := <-m.addChan:
				m.data[metric.Name] = append(
					m.data[metric.Name],
					MetricBase.MetricValues{Time: metric.Time, Value: metric.Value},
				)
			case req := <-m.listRequests:
				for key := range m.data {
					req <- key
				}
				close(req)
			case req := <-m.dataRequests:
				if _, ok := m.data[req.Name]; ok {
					for _, data := range m.data[req.Name] {
						req.Result <- data
					}
				}
				close(req.Result)
			case <-m.stopChan:
				close(m.stopChan)
				close(m.addChan)
				close(m.listRequests)
				close(m.dataRequests)
				return
			}
		}
	}()
}

func (m *MemoryBackend) Stop() { m.stopChan <- true }

func (m *MemoryBackend) AddMetric(metric MetricBase.Metric) {
	m.addChan <- metric
}

func (m *MemoryBackend) GetMetricsList(results chan string) {
	m.listRequests <- results
}

func (m *MemoryBackend) GetRawData(name string, from, to int64, result chan MetricBase.MetricValues) {
	m.dataRequests <- dataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}
