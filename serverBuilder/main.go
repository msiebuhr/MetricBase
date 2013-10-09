package serverBuilder

import (
	"github.com/msiebuhr/MetricBase"
)

type MetricServer struct {
	frontends []MetricBase.Frontend
	backends  []MetricBase.Backend
	stopChan  chan bool
}

func (m *MetricServer) AddFrontend(f MetricBase.Frontend) {
	m.frontends = append(m.frontends, f)

	// Hook up to the first backend
	if len(m.backends) >= 1 {
		f.SetBackend(m.backends[0])
	}
}

func (m *MetricServer) AddBackend(b MetricBase.Backend) {
	m.backends = append(m.backends, b)

	// Make the last backend point to this one
	if len(m.backends) > 1 {
		m.backends[len(m.backends)-2].SetBackend(b)
	}

	// If we insert the first backend, hook up any frontend to this
	if len(m.backends) == 1 {
		for i := range m.frontends {
			m.frontends[i].SetBackend(b)
		}
	}
}

func (m *MetricServer) Start() {
	// Start all back-ends in reverse order
	for i := len(m.backends) - 1; i >= 0; i-- {
		go m.backends[i].Start()
	}

	// Start all front-ends, now they can talk to something
	for i := range m.frontends {
		go m.frontends[i].Start()
	}

	// Wait for order to stop
	<-m.stopChan

	// Close up front-ends
	for i := range m.frontends {
		m.frontends[i].Stop()
	}

	for i := len(m.backends) - 1; i >= 0; i-- {
		m.backends[i].Stop()
	}
}

func (m *MetricServer) Stop() {
	m.stopChan <- true
}

func CreateMetricServer() MetricServer {
	return MetricServer{stopChan: make(chan bool)}
}
