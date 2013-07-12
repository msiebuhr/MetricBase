package frontends

import (
	"encoding/json"
	"fmt"
	"github.com/msiebuhr/MetricBase"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type HttpServer struct {
	staticRoot string
	backend    MetricBase.Backend
}

func CreateHttpServer(staticRoot string) *HttpServer {
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
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (h *HttpServer) Stop() {
	// NOP - no way of stopping a HTTP server, aparently
}
