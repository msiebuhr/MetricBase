package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/frontends"
	"github.com/msiebuhr/MetricBase/serverBuilder"
)

var staticRoot = flag.String("http-pub", "./http-pub", "HTTP public dir")
var boltDb = flag.String("boltdb", "./bolt.db", "Bolt db file")

func main() {
	// Parse command line flags
	flag.Parse()

	// Create backend + database
	bdb, err := backends.NewBoltBackend(*boltDb)
	if err != nil {
		fmt.Println("Could not create bolt database", err)
		return
	}

	// Create server
	mb := serverBuilder.NewMetricServer(
		[]frontends.Frontend{
			frontends.NewHttpServer(*staticRoot),
			frontends.NewGraphiteTcpServer(),
		},
		//backends.NewTestProxy(backends.NewMemoryBackend()),
		backends.NewTestProxy(bdb),
	)

	go mb.Start()

	// Listen for signals and stop
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Stopping server:", <-ch)
	mb.Stop()
}
