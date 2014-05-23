package frontends

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/msiebuhr/MetricBase/backends"
	"github.com/msiebuhr/MetricBase/metrics"
	"github.com/msiebuhr/MetricBase/query"
)

type HttpServer struct {
	staticRoot string
	backend    backends.Backend
}

func NewHttpServer(staticRoot string) *HttpServer {
	absRoot, err := filepath.Abs(staticRoot)
	if err != nil {
		absRoot = staticRoot
	}
	return &HttpServer{staticRoot: absRoot}
}

func (h *HttpServer) SetBackend(backend backends.Backend) {
	h.backend = backend
}

func (h *HttpServer) listHandler(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(backends.GetMetricsAsList(h.backend))
	if err != nil {
		http.Error(w, "Could not Encode JSON", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (h *HttpServer) metricHandler(w http.ResponseWriter, req *http.Request) {
	// Parse the URL
	urlParts := strings.Split(req.URL.Path[1:], "/")
	if len(urlParts) != 3 {
		http.NotFound(w, req)
		return
	}

	// New request
	resultChan := make(chan metrics.MetricValue, 100)
	h.backend.GetRawData(urlParts[2], time.Now().Add(-1*time.Hour*24*7), time.Now(), resultChan)

	// Fetch the data
	newData := make(map[string]float64)
	for data := range resultChan {
		newData[fmt.Sprintf("%v", data.Time.Unix())] = data.Value
	}

	// Encode as JSON
	b, err := json.Marshal(newData)
	if err != nil {
		http.Error(w, "Could not Encode JSON", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (h *HttpServer) queryHandler(w http.ResponseWriter, req *http.Request) {
	// Populate query
	req.ParseForm()

	// Parse start and end times
	startDateString := req.FormValue("start")
	endDateString := req.FormValue("end")

	if startDateString == "" {
		startDateString = "-1w"
	}
	if endDateString == "" {
		endDateString = "0"
	}

	// Parse the dates
	// TODO: Know about absolute dates as well
	startDuration, err := parseDuration(startDateString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	endDuration, err := parseDuration(endDateString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Switch, so we start with the earliest date
	if startDuration > endDuration {
		endDuration, startDuration = startDuration, endDuration
	}

	// Build the query
	res, err := query.ParseGraphiteQuery(req.FormValue("q"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ARGH. Tends to hang about here somewhere...
	responses, err := res.Query(query.Request{
		Backend: h.backend,
		From:    time.Now().Add(startDuration),
		To:      time.Now().Add(endDuration),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the data
	newData := make(map[string]map[string]float64)
	for _, result := range responses {
		name := result.Meta["name"]
		newData[name] = make(map[string]float64)
		for data := range result.Data {
			//fmt.Println("Got data", data)
			newData[name][fmt.Sprintf("%v", data.Time.Unix())] = data.Value
		}
	}

	// Encode as JSON
	b, err := json.Marshal(newData)
	if err != nil {
		http.Error(w, "Could not Encode JSON", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (h *HttpServer) staicHandler(w http.ResponseWriter, req *http.Request) {
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
	http.HandleFunc("/rpc/list", h.listHandler)
	http.HandleFunc("/rpc/get/", h.metricHandler)
	http.HandleFunc("/rpc/query", h.queryHandler)
	http.HandleFunc("/", h.staicHandler)
	fmt.Println("Web interface on http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Stop is a No-op, as Go's http server doesn't support being stopped.
func (h *HttpServer) Stop() {
	// NOP - no way of stopping a HTTP server, aparently
}
