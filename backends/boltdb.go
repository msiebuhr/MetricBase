package backends

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/msiebuhr/MetricBase/metrics"
)

type BoltBackend struct {
	db *bolt.DB

	stopChan     chan bool
	addChan      chan metrics.Metric
	listRequests chan chan string
	dataRequests chan dataRequest
}

func NewBoltBackend(filename string) (*BoltBackend, error) {
	db, err := bolt.Open(filename, 0666)
	if err != nil {
		return nil, err
	}

	return &BoltBackend{
		db:           db,
		stopChan:     make(chan bool),
		addChan:      make(chan metrics.Metric),
		listRequests: make(chan chan string),
		dataRequests: make(chan dataRequest),
	}, nil
}

func serializeMetric(m metrics.Metric) (key []byte, value []byte) {
	// Encode key
	key = []byte(fmt.Sprintf("%013d", m.Time))

	// Encode value
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, m.Value)
	value = buf.Bytes()

	return key, value
}

// {{{
// Serialize metrics
func parseMetricKey(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 10, 64)
}

func parseMetricValue(b []byte) float64 {
	var r float64
	buf := bytes.NewBuffer(b)
	_ = binary.Read(buf, binary.LittleEndian, &r)
	return r
}

func (m *BoltBackend) Start() {
	go func() {
		for {
			select {
			case metric := <-m.addChan:
				// Insert it in a database transaction
				err := m.db.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists([]byte(metric.Name))
					if err != nil {
						return err
					}

					key, value := serializeMetric(metric)

					err = b.Put(key, value)
					if err != nil {
						return err
					}

					return nil
				})

				if err != nil {
					fmt.Errorf("Bolt was unhappy writing data: %v", err)
				}

			case req := <-m.listRequests:
				// List all buckets in a view-transaction
				m.db.View(func(tx *bolt.Tx) error {
					c := tx.Cursor()
					k, _ := c.First()
					for k != nil {
						req <- string(k)
						k, _ = c.Next()
					}
					return nil
				})
				close(req)

			case req := <-m.dataRequests:
				// List all relevant data in a view-transaction
				m.db.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte(req.Name))

					// Bail if bucket does not exist
					if b == nil {
						fmt.Println("Bucket does not exist")
						return nil
					}

					// It's kosher - return relevant data
					cursor := b.Cursor()

					// Generate the first key, seek to it and begin reading data
					firstKey := []byte(fmt.Sprintf("%013d", req.From))
					rawTime, rawVal := cursor.Seek(firstKey)
					time, _ := parseMetricKey(rawTime)

					// Loop over the rest
					for time <= req.To {
						// Send the data
						req.Result <- metrics.MetricValue{
							Time:  time,
							Value: parseMetricValue(rawVal),
						}

						// Extract a new time/value
						rawTime, rawVal = cursor.Next()

						// Break from the loop if we reach the end
						if rawTime == nil {
							break
						}

						// Parse time and loop-de-loop
						time, _ = parseMetricKey(rawTime)
					}

					return nil
				})

				close(req.Result)

			case <-m.stopChan:
				// Stop comm chans and shut down database
				close(m.stopChan)
				close(m.addChan)
				close(m.listRequests)
				close(m.dataRequests)
				m.db.Close()
				return
			}
		}
	}()
}

func (m *BoltBackend) Stop() { m.stopChan <- true }

func (m *BoltBackend) AddMetric(metric metrics.Metric) {
	m.addChan <- metric
}

func (m *BoltBackend) GetMetricsList(results chan string) {
	m.listRequests <- results
}

func (m *BoltBackend) GetRawData(name string, from, to int64, result chan metrics.MetricValue) {
	m.dataRequests <- dataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}
