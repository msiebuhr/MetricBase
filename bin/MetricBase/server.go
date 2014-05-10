package main

import (
	"fmt"
	"github.com/msiebuhr/MetricBase"
	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/frontends"
	"github.com/msiebuhr/MetricBase/serverBuilder"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create server
	mb := serverBuilder.NewMetricServer(
		[]MetricBase.Frontend{
			frontends.NewHttpServer("./http-pub"),
			frontends.NewGraphiteTcpServer(),
		},
		backends.NewTestProxy(backends.NewMemoryBackend()),
	)

	go mb.Start()

	// Listen for signals and stop
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Stopping server:", <-ch)
	mb.Stop()
}
