package backends

import (
	"time"

	"github.com/msiebuhr/MetricBase/metrics"
)

type MemoryBackend struct {
	data map[string][]metrics.MetricValue

	stopChan     chan bool
	addChan      chan metrics.Metric
	listRequests chan chan string
	dataRequests chan dataRequest
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		data:         make(map[string][]metrics.MetricValue),
		stopChan:     make(chan bool),
		addChan:      make(chan metrics.Metric),
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
					metric.GetMetricValue(),
				)
			case req := <-m.listRequests:
				for key := range m.data {
					req <- key
				}
				close(req)
			case req := <-m.dataRequests:
				if _, ok := m.data[req.Name]; ok {
					for _, data := range m.data[req.Name] {
						if data.Time.After(req.From) && data.Time.Before(req.To) {
							req.Result <- data
						}
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

func (m *MemoryBackend) AddMetric(metric metrics.Metric) {
	m.addChan <- metric
}

func (m *MemoryBackend) GetMetricsList(results chan string) {
	m.listRequests <- results
}

func (m *MemoryBackend) GetRawData(name string, from, to time.Time, result chan metrics.MetricValue) {
	m.dataRequests <- dataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}
