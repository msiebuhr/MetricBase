// In-memory metric store
package backends

import (
	"github.com/msiebuhr/MetricBase"
)

type MemoryBackend struct {
	data map[string][]MetricBase.MetricValues

	stopChan        chan bool
	addChan         chan MetricBase.Metric
	listRequestChan chan MetricBase.ListRequest
	dataRequestChan chan MetricBase.DataRequest
}

func CreateMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		data:            make(map[string][]MetricBase.MetricValues),
		stopChan:        make(chan bool),
		addChan:         make(chan MetricBase.Metric),
		listRequestChan: make(chan MetricBase.ListRequest),
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
					req.Result <- key
				}
				close(req.Result)
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
func (m *MemoryBackend) Add(metrics MetricBase.AddRequest) {
	go func() {
		for metric := range metrics.Data {
			m.addChan <- metric
		}
	}()
}
func (m *MemoryBackend) List(req MetricBase.ListRequest) {
	m.listRequestChan <- req
}
func (m *MemoryBackend) Data(req MetricBase.DataRequest) {
	m.dataRequestChan <- req
}

func (m *MemoryBackend) SetBackend(backend MetricBase.Backend) {
	// NOP - the buck stops here
}
