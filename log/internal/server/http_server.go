package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHttpServer(addr *string) *http.Server {
	httpserver := newHttpServer(addr)
	router := mux.NewRouter()
	router.HandleFunc("/", httpserver.handleProduce).Methods("POST")
	router.HandleFunc("/", httpserver.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    *addr,
		Handler: router,
	}
}

type httpServer struct {
	Log *Log
}

func newHttpServer(addr *string) *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}
type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (server httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := server.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ProduceResponse{Offset: offset}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (server httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := server.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
