package main

import (
	"github.com/msiebuhr/MetricBase"
	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/frontends"
)

func main() {
	// Create server
	mb := MetricBase.CreateMetricBaseServer()

	// Create and add front- and back-ends

	mb.AddFrontend(frontends.CreateHttpServer("./"))

	mb.AddBackend(backends.CreateMemoryBackend())

	mb.Start()
}
