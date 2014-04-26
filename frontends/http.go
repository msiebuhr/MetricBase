package frontends

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/msiebuhr/MetricBase"
	"github.com/msiebuhr/MetricBase/backends"
)

type HttpServer struct {
	staticRoot string
	backend    MetricBase.Backend
}

func NewHttpServer(staticRoot string) *HttpServer {
	absRoot, err := filepath.Abs(staticRoot)
	if err != nil {
		absRoot = staticRoot
	}
	return &HttpServer{staticRoot: absRoot}
}

func (h *HttpServer) SetBackend(backend MetricBase.Backend) {
	h.backend = backend
}

func (h *HttpServer) GetList(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(backends.GetMetricsAsList(h.backend))
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
	resultChan := make(chan MetricBase.MetricValues, 100)
	h.backend.GetRawData(urlParts[2], 0, 0, resultChan)

	// Fetch the data
	newData := make(map[string]float64)
	for data := range resultChan {
		newData[fmt.Sprintf("%v", data.Time)] = data.Value
	}

	// Encode as JSON
	b, err := json.Marshal(newData)
	if err != nil {
		http.Error(w, "Could not Encode JSON", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (h *HttpServer) GetStatic(w http.ResponseWriter, req *http.Request) {
	// Figure out the path
	abspath, err := filepath.Abs(filepath.Join(h.staticRoot, req.URL.Path[1:]))
	if err != nil {
		http.Error(w, "Could not figure out path", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(abspath, h.staticRoot) {
		http.Error(w, "Invalid path", http.StatusInternalServerError)
		return
	}

	// Takes care of checking the file exists, hunt down an index.html if it's
	// a directory, setting correct content-type, range-requests, ...
	http.ServeFile(w, req, abspath)
}

func (h *HttpServer) Start() {
	http.HandleFunc("/rpc/list", h.GetList)
	http.HandleFunc("/rpc/get/", h.GetMetric)
	http.HandleFunc("/", h.GetStatic)
	fmt.Println("Web interface on http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (h *HttpServer) Stop() {
	// NOP - no way of stopping a HTTP server, aparently
}
