package frontends

import (
	"bufio"
	"fmt"
	"github.com/msiebuhr/MetricBase"
	"net"
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

func (g *GraphiteTcpServer) handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	defer conn.Close()

	// Create addition-channel to backend
	addChan := make(chan MetricBase.Metric, 100)
	g.backend.AddMetrics(addChan)
	defer close(addChan)

	for scanner.Scan() {
		// PARSE METRIC LINES
		//err, m :=
		err, m := MetricBase.ParseGraphiteLine(scanner.Text())
		if err != nil {
			// Panic / close connection?
			break
		}
		// Send parsed metric to the back-end
		addChan <- m
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error while parsing text: %v", err)
	}
}

func (g *GraphiteTcpServer) Start() {
	server, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Listening on localhost:8000")

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
