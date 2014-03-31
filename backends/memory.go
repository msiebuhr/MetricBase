// In-memory metric store
package backends

import (
	"github.com/msiebuhr/MetricBase"
)

type MemoryBackend struct {
	data map[string][]MetricBase.MetricValues

	stopChan        chan bool
	addChan         chan MetricBase.Metric
	listRequestChan chan chan string
	dataRequestChan chan MetricBase.DataRequest
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		data:            make(map[string][]MetricBase.MetricValues),
		stopChan:        make(chan bool),
		addChan:         make(chan MetricBase.Metric),
		listRequestChan: make(chan chan string),
		dataRequestChan: make(chan MetricBase.DataRequest),
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
			case req := <-m.listRequestChan:
				for key := range m.data {
					req <- key
				}
				close(req)
			case req := <-m.dataRequestChan:
				if _, ok := m.data[req.Name]; ok {
					for _, data := range m.data[req.Name] {
						req.Result <- data
					}
				}
				close(req.Result)
			case <-m.stopChan:
				close(m.stopChan)
				close(m.addChan)
				close(m.listRequestChan)
				close(m.dataRequestChan)
				return
			}
		}
	}()
}

func (m *MemoryBackend) Stop() { m.stopChan <- true }

func (m *MemoryBackend) AddMetrics(metrics chan MetricBase.Metric) {
	go func() {
		for metric := range metrics {
			m.addChan <- metric
		}
	}()
}

func (m *MemoryBackend) GetMetricsList(results chan string) {
	m.listRequestChan <- results
}

func (m *MemoryBackend) GetRawData(name string, from, to int64, result chan MetricBase.MetricValues) {
	m.dataRequestChan <- MetricBase.DataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}

func (m *MemoryBackend) SetBackend(backend MetricBase.Backend) {
	// NOP - the buck stops here
}
