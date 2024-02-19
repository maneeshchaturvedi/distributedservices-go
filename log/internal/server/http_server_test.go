package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHandleProduce(t *testing.T) {
	log := NewLog()
	server := httpServer{Log: log}

	record := Record{
		Value:  []byte("test value"),
		Offset: 0,
	}
	req := ProduceRequest{Record: record}
	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	server.handleProduce(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("handleProduce() returned incorrect status code, expected %d but got %d", http.StatusOK, w.Code)
	}

	var res ProduceResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	if err != nil {
		t.Errorf("handleProduce() returned an error while decoding response body: %v", err)
	}

	if res.Offset != 0 {
		t.Errorf("handleProduce() returned incorrect offset, expected 0 but got %d", res.Offset)
	}

	if len(log.records) != 1 {
		t.Errorf("handleProduce() did not add the record to the log")
	}
}

func TestHandleConsume(t *testing.T) {
	log := NewLog()
	server := httpServer{Log: log}

	record := Record{
		Value:  []byte("test value"),
		Offset: 0,
	}
	log.records = append(log.records, record)

	req := ConsumeRequest{Offset: 0}
	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))

	server.handleConsume(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("handleConsume() returned incorrect status code, expected %d but got %d", http.StatusOK, w.Code)
	}

	var res ConsumeResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	if err != nil {
		t.Errorf("handleConsume() returned an error while decoding response body: %v", err)
	}

	if !cmp.Equal(res.Record, record) {
		t.Errorf("handleConsume() returned incorrect record, expected %+v but got %+v", record, res.Record)
	}
}

func TestHandleConsume_NonExistentOffset(t *testing.T) {
	log := NewLog()
	server := httpServer{Log: log}

	req := ConsumeRequest{Offset: 0}
	body, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))

	server.handleConsume(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("handleConsume() returned incorrect status code, expected %d but got %d", http.StatusNotFound, w.Code)
	}

}
