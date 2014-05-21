package backends

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/msiebuhr/MetricBase/metrics"
)

type BoltBackend struct {
	db            *bolt.DB
	addBuffer     map[string][]metrics.Metric
	addBufferSize uint32

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

	// We generally do a lot of appending, so setting FillPercent high makes a
	// lot of sense (BoltDB docs + simple benchmarks give up to 10% speedup)
	db.FillPercent = 95

	return &BoltBackend{
		db:            db,
		addBuffer:     make(map[string][]metrics.Metric),
		addBufferSize: 0,

		stopChan:     make(chan bool),
		addChan:      make(chan metrics.Metric),
		listRequests: make(chan chan string),
		dataRequests: make(chan dataRequest),
	}, nil
}

func putFloat64(v float64) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, v)
	return buf.Bytes()
}

func putUint40(v uint64) []byte {
	b := make([]byte, 5)
	b[0] = byte(v >> 32)
	b[1] = byte(v >> 24)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 8)
	b[4] = byte(v)
	return b
}

func parseUint40(b []byte) uint64 {
	return uint64(b[4]) | uint64(b[3])<<8 | uint64(b[2])<<16 | uint64(b[1])<<24 |
		uint64(b[0])<<32
}

func parseFloat64(b []byte) float64 {
	var r float64
	buf := bytes.NewBuffer(b)
	_ = binary.Read(buf, binary.LittleEndian, &r)
	return r
}

func (m *BoltBackend) flushAddBuffer() {
	// Insert it in a database transaction
	err := m.db.Update(func(tx *bolt.Tx) error {
		for seriesName, seriesData := range m.addBuffer {
			b, err := tx.CreateBucketIfNotExists([]byte(seriesName))
			if err != nil {
				return err
			}

			for _, m := range seriesData {
				err = b.Put(putUint40(uint64(m.Time.Unix())), putFloat64(m.Value))
				if err != nil {
					return err
				}
			}

			// Remove data from buffer
			delete(m.addBuffer, seriesName)
		}
		return nil
	})

	m.addBuffer = make(map[string][]metrics.Metric)
	m.addBufferSize = 0

	if err != nil {
		fmt.Errorf("Bolt was unhappy writing data: %v", err)
	}

}

func (m *BoltBackend) Start() {
	go func() {
		t := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-t.C:
				// Regularly flush the internal buffer
				m.flushAddBuffer()

			case metric := <-m.addChan:
				// Add it to the internal buffer
				m.addBuffer[metric.Name] = append(m.addBuffer[metric.Name], metric)
				m.addBufferSize += 1

				// Break if our buffer is too small
				if m.addBufferSize < 10000 {
					break
				}

				m.flushAddBuffer()

			case req := <-m.listRequests:
				// List all buckets in a view-transaction
				m.db.View(func(tx *bolt.Tx) error {
					c := tx.Cursor()
					k, _ := c.First()
					for k != nil {
						// Skip sending this metric if it's in the buffer
						if _, ok := m.addBuffer[string(k)]; !ok {
							req <- string(k)
						}
						k, _ = c.Next()
					}
					return nil
				})

				// Add whatever we find in the buffer
				for name := range m.addBuffer {
					req <- name
				}

				close(req)

			case req := <-m.dataRequests:
				// List all relevant data in a view-transaction
				m.db.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte(req.Name))

					// Bail if bucket does not exist
					if b == nil {
						fmt.Printf("Bolt - no series '%v'.\n", req.Name)
						return nil
					}

					// It's kosher - return relevant data
					cursor := b.Cursor()

					// Generate the first key, seek to it and begin reading data
					firstKey := putUint40(uint64(req.From.Unix()))
					rawTime, rawVal := cursor.Seek(firstKey)
					metricTime := time.Unix(int64(parseUint40(rawTime)), 0)

					// Loop over the rest
					for metricTime.Before(req.To) {
						// Send the data
						req.Result <- metrics.MetricValue{
							Time:  metricTime,
							Value: parseFloat64(rawVal),
						}

						// Extract a new time/value
						rawTime, rawVal = cursor.Next()

						// Break from the loop if we reach the end
						if rawTime == nil {
							break
						}

						// Parse time and loop-de-loop
						metricTime = time.Unix(int64(parseUint40(rawTime)), 0)
					}

					return nil
				})

				// Add stuff from the buffer
				if data, ok := m.addBuffer[req.Name]; ok {
					for _, m := range data {
						if m.Time.After(req.From) && m.Time.Before(req.To) {
							req.Result <- metrics.MetricValue{
								Time:  m.Time,
								Value: m.Value,
							}
						}
					}
				}

				close(req.Result)

			case <-m.stopChan:
				// Stop comm chans and shut down database
				t.Stop()
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

func (m *BoltBackend) GetRawData(name string, from, to time.Time, result chan metrics.MetricValue) {
	m.dataRequests <- dataRequest{
		Name:   name,
		From:   from,
		To:     to,
		Result: result,
	}
}
