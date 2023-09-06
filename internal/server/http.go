package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Log *Log
}

type ProduceRequest struct {
	Record Record `json:"record"`
}
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}
type ConsumeResponse struct {
	Record Record `json:"record"`
}

func newHttpServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func NewHttpServer(addr string) *http.Server {
	httpsrv := newHttpServer()
	router := mux.NewRouter()

	router.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	router.HandleFunc("/", httpsrv.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func (serverLog *httpServer) handleProduce(writer http.ResponseWriter, request *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(request.Body).Decode(&req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := serverLog.Log.Append(req.Record)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ProduceResponse{Offset: offset}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (serverLog *httpServer) handleConsume(writer http.ResponseWriter, request *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(request.Body).Decode(&req)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := serverLog.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ConsumeResponse{Record: record}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
