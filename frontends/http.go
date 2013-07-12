package frontends

import (
	"encoding/json"
	"fmt"
	"github.com/msiebuhr/MetricBase"
	"log"
	"net/http"
	"strings"
)

type HttpServer struct {
	staticRoot string
	backend    MetricBase.Backend
}

func CreateHttpServer(staticRoot string) *HttpServer {
	return &HttpServer{staticRoot: staticRoot}
}

func (h *HttpServer) SetBackend(backend MetricBase.Backend) {
	h.backend = backend
}

func (h *HttpServer) GetList(w http.ResponseWriter, req *http.Request) {
	listReq := &MetricBase.ListRequest{Result: make(chan string, 10)}
	listRes := make([]string, 0)
	h.backend.List(*listReq)

	for res := range listReq.Result {
		listRes = append(listRes, res)
	}

	b, err := json.Marshal(listRes)
	if err != nil {
		http.Error(w, "Could not Encode JSON", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (h *HttpServer) GetMetric(w http.ResponseWriter, req *http.Request) {
	// Parse the URL
	urlParts := strings.Split(req.URL.Path[1:], "/")
	if len(urlParts) != 3 {
		http.NotFound(w, req)
		return
	}

	// New request
	dataReq := &MetricBase.DataRequest{
		Name:   urlParts[2],
		Result: make(chan MetricBase.MetricValues),
	}

	// Fetch the data
	h.backend.Data(*dataReq)
	newData := make(map[string]float64)
	for data := range dataReq.Result {
		newData[fmt.Sprintf("%v", data.Time)] = data.Value
	}

	// Encode as JSON
	b, err := json.Marshal(newData)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	w.Write(b)
}

func (h *HttpServer) GetStatic(w http.ResponseWriter, req *http.Request) {
	// Return whatever static file we find...
	fmt.Fprintf(w, "Serve static file %v.", req.URL.Path[1:])
}

func (h *HttpServer) Start() {
	http.HandleFunc("/rpc/list", h.GetList)
	http.HandleFunc("/rpc/get/", h.GetMetric)
	//http.HandleFunc("/", h.GetStatic)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (h *HttpServer) Stop() {
	// NOP - no way of stopping a HTTP server, aparently
}
