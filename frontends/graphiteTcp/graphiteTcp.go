package graphiteTcp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/metrics"
)

type GraphiteTcpServer struct {
	backend backends.Backend
}

func NewGraphiteTcpServer() *GraphiteTcpServer {
	return &GraphiteTcpServer{}
}

func (g *GraphiteTcpServer) SetBackend(backend backends.Backend) {
	g.backend = backend
}

func (g *GraphiteTcpServer) handleConnection(conn io.ReadWriteCloser) {
	scanner := bufio.NewScanner(conn)
	defer conn.Close()

	// Create a channel for the metrics
	addChan := make(chan metrics.Metric, 1000)
	g.backend.AddMetricChan(addChan)
	defer close(addChan)

	for scanner.Scan() {
		// PARSE METRIC LINES
		//err, m :=
		err, m := parseGraphiteLine(scanner.Text())
		if err == nil {
			// Send parsed metric to the back-end
			addChan <- m
		} else {
			conn.Write([]byte(err.Error()))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error while parsing text: %v", err)
	}
}

func parseGraphiteLine(line string) (error, metrics.Metric) {
	// Find newline-rune
	fields := strings.Fields(line)

	newMetric := metrics.Metric{}

	// A line must at least contain <metric.name> <timestamp> <value> <tag=value>+
	if len(fields) != 3 {
		return errors.New("Invalid line"), newMetric
	}

	// Convert name
	newMetric.Name = string(fields[0])

	// Parse out value
	value, err := strconv.ParseFloat(string(fields[1]), 64)
	if err != nil {
		return err, newMetric
	}
	newMetric.Value = value

	// Parse out timestamp
	timeNum, err := strconv.ParseInt(string(fields[2]), 10, 64)
	if err != nil {
		return err, newMetric
	}
	newMetric.Time = time.Unix(timeNum, 0)

	return nil, newMetric
}

func (g *GraphiteTcpServer) Start() {
	server, err := net.Listen("tcp", ":2003")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Graphite TCP interface on :2003")

	// Listen for connections
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go g.handleConnection(conn)
	}
}

func (g *GraphiteTcpServer) Stop() {
}
