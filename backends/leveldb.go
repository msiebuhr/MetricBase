package backends

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/jmhodges/levigo"
	"github.com/msiebuhr/MetricBase"
	"strconv"
	"strings"
)

func serializeMetric(m MetricBase.Metric) (key []byte, value []byte) {
	// Encode key
	key = []byte(fmt.Sprintf("%v:%013d", m.Name, m.Time))

	// Encode value
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, m.Value)
	value = buf.Bytes()

	return key, value
}

// {{{
// Serialize metrics
func parseMetricKey(b []byte) (error, string, int64) {
	parts := strings.SplitN(string(b), ":", 2)

	name := parts[0]
	time, err := strconv.ParseInt(parts[1], 10, 64)

	return err, name, time
}

func parseMetricValue(b []byte) float64 {
	var r float64
	buf := bytes.NewBuffer(b)
	_ = binary.Read(buf, binary.LittleEndian, &r)
	return r
}

// }}}

type LevelDb struct {
	store *levigo.DB

	addRequests  chan MetricBase.Metric
	listRequests chan chan string
	dataRequests chan MetricBase.DataRequest

	stopChan chan bool
}

func NewLevelDb(filename string) *LevelDb {
	options := levigo.NewOptions()
	options.SetCreateIfMissing(true)

	db, err := levigo.Open(filename, options)

	if err != nil {
		panic(err)
	}

	ls := &LevelDb{
		addRequests:  make(chan MetricBase.Metric, 100),
		listRequests: make(chan chan string, 10),
		dataRequests: make(chan MetricBase.DataRequest, 10),
		stopChan:     make(chan bool),
		store:        db,
	}

	return ls
}

func (l *LevelDb) Start() {
	// Start listener-loop
	go func() {
		for {
			select {
			case <-l.stopChan:
				l.store.Close()
				return
			case metric := <-l.addRequests:
				l.addMetric(metric)
			case query := <-l.listRequests:
				l.listMetrics(query)
			case query := <-l.dataRequests:
				l.handleData(query)
			}
		}
	}()
}

func (l *LevelDb) AddMetrics(metrics chan MetricBase.Metric) {
	go func() {
		for m := range metrics {
			l.addRequests <- m
		}
	}()
}

func (l *LevelDb) Stop() {
	l.stopChan <- true
}

func (l *LevelDb) addMetric(metric MetricBase.Metric) {
	wo := levigo.NewWriteOptions()
	defer wo.Close()
	k, v := serializeMetric(metric)
	_ = l.store.Put(wo, k, v)
}

func (l *LevelDb) listMetrics(result chan string) {
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	iter := l.store.NewIterator(ro)
	defer iter.Close()

	iter.Seek([]byte{0x00})

	var currentName string = ""

	for iter = iter; iter.Valid(); iter.Next() {
		err, name, _ := parseMetricKey(iter.Key())

		// Ignore errors
		if err != nil {
			continue
		}

		// Ignore similar names.
		if name != currentName {
			result <- name
			currentName = name
		}
	}

	close(result)
}

func (l *LevelDb) handleData(query MetricBase.DataRequest) {
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	iter := l.store.NewIterator(ro)
	defer iter.Close()

	iter.Seek([]byte(fmt.Sprintf("%v:", query.Name)))
	for iter = iter; iter.Valid(); iter.Next() {
		err, name, time := parseMetricKey(iter.Key())
		value := parseMetricValue(iter.Value())

		if name != query.Name {
			break
		}

		if err != nil {
			continue
		}

		query.Result <- MetricBase.MetricValues{Time: time, Value: value}
	}
	close(query.Result)
}

func (l *LevelDb) GetMetricsList(results chan string) {
	l.listRequests <- results
}

func (l *LevelDb) GetRawData(name string, from, to int64, result chan MetricBase.MetricValues) {
	l.dataRequests <- MetricBase.DataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}
