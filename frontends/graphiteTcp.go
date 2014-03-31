package frontends

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/msiebuhr/MetricBase"
	"io"
	"net"
	"strconv"
	"strings"
)

type GraphiteTcpServer struct {
	backend MetricBase.Backend
}

func CreateGraphiteTcpServer() *GraphiteTcpServer {
	return &GraphiteTcpServer{}
}

func (g *GraphiteTcpServer) SetBackend(backend MetricBase.Backend) {
	g.backend = backend
}

func (g *GraphiteTcpServer) handleConnection(conn io.ReadWriteCloser) {
	scanner := bufio.NewScanner(conn)
	defer conn.Close()

	// Create addition-channel to backend
	addChan := make(chan MetricBase.Metric, 100)
	g.backend.AddMetrics(addChan)
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
		fmt.Println("Error while parsing text: %v", err)
	}
}

func parseGraphiteLine(line string) (error, MetricBase.Metric) {
	// Find newline-rune
	fields := strings.Fields(line)

	newMetric := MetricBase.Metric{}

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
	time, err := strconv.ParseInt(string(fields[2]), 10, 64)
	if err != nil {
		return err, newMetric
	}
	newMetric.Time = time

	return nil, newMetric
}

func (g *GraphiteTcpServer) Start() {
	server, err := net.Listen("tcp", ":2003")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Listening on localhost:2003")

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
