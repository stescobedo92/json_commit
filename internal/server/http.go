package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stescobedo92/json_commit/internal/log"
)

type httpServer struct {
	Log *log.Log
}

func NewHTTPServer(addr string, log *log.Log) *http.Server {
	srv := &httpServer{
		Log: log,
	}
	r := mux.NewRouter()

	r.HandleFunc("/records", srv.handleProduce).Methods("POST")
	r.HandleFunc("/records/{offset}", srv.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type ProduceRequest struct {
	Record log.Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record log.Record `json:"record"`
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	offset, err := s.Log.Append(req.Record)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, ProduceResponse{Offset: offset})
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	offsetStr, ok := vars["offset"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "'offset' parameter is missing in URL")
		return
	}

	offset, err := strconv.ParseUint(offsetStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid 'offset' parameter")
		return
	}

	record, err := s.Log.Read(offset)
	if err != nil {
		if err == log.ErrOffsetNotFound {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ConsumeResponse{Record: record})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error: failed to marshal JSON response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
