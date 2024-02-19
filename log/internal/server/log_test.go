package server

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLogAppend(t *testing.T) {
	log := NewLog()

	record := Record{
		Value:  []byte("test value"),
		Offset: 0,
	}

	offset, err := log.Append(record)
	if err != nil {
		t.Errorf("Append() returned an error: %v", err)
	}

	if offset != 0 {
		t.Errorf("Append() returned incorrect offset, expected 0 but got %d", offset)
	}

	if len(log.records) != 1 {
		t.Errorf("Append() did not add the record to the log")
	}
}

func TestLogRead(t *testing.T) {
	log := NewLog()

	record := Record{
		Value:  []byte("test value"),
		Offset: 0,
	}

	log.records = append(log.records, record)

	readRecord, err := log.Read(0)
	if err != nil {
		t.Errorf("Read() returned an error: %v", err)
	}

	if !cmp.Equal(readRecord, record) {
		t.Errorf("Read() returned incorrect record, expected %+v but got %+v", record, readRecord)
	}

	_, err = log.Read(1)
	if err == nil {
		t.Errorf("Read() did not return an error for non-existent offset")
	}
}
