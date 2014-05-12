package serverBuilder

import (
	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/frontends"
)

type MetricServer struct {
	frontends []frontends.Frontend
	backend   backends.Backend
	stopChan  chan bool
}

// NewMetricServer returns a new MetricServer with the given backend and list
// of frontends.
func NewMetricServer(f []frontends.Frontend, b backends.Backend) MetricServer {
	// Hook up backends
	for _, front := range f {
		front.SetBackend(b)
	}

	// Server construction
	return MetricServer{
		stopChan:  make(chan bool),
		frontends: f,
		backend:   b,
	}
}

// Start the server. No guarantees are made about re-startability.
func (m *MetricServer) Start() {
	// Start the backend
	go m.backend.Start()

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

	m.backend.Stop()
}

// Stop the server
func (m *MetricServer) Stop() {
	m.stopChan <- true
}
