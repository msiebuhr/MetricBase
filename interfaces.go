package MetricBase

// Data structures
type MetricValues struct {
	Time  int64
	Value float64
}

type Metric struct {
	MetricValues
	Name string
}

// Requests
type AddRequest struct {
	Data chan Metric
}

type ListRequest struct {
	Result chan string
}

type DataRequest struct {
	Name   string
	From   int64
	To     int64
	Result chan MetricValues
}

// TODO: Helpers to generate these requests

// Interfaces
type Backend interface {
	AddBackend(Backend)
	Start()
	Stop()
	Add(AddRequest)
	List(ListRequest)
	Data(DataRequest)
}

type Frontend interface {
	AddBackend(Backend)
	Start()
	Stop()
}
