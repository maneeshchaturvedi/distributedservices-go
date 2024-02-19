package server

import (
	"fmt"
	"sync"
)

// Log represents a log that stores records.
type Log struct {
	mu      sync.Mutex
	records []Record
}

// NewLog creates a new instance of Log.
func NewLog() *Log {
	return &Log{}
}

// Append appends a record to the log and returns its offset.
func (log *Log) Append(record Record) (uint64, error) {
	log.mu.Lock()
	defer log.mu.Unlock()
	record.Offset = uint64(len(log.records))
	log.records = append(log.records, record)
	return record.Offset, nil
}

// Read retrieves a record from the log based on the given offset.
func (log *Log) Read(offset uint64) (Record, error) {
	log.mu.Lock()
	defer log.mu.Unlock()
	if offset >= uint64(len(log.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return log.records[offset], nil
}

// Record represents a log record.
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")
