package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Fatalf("Internal server error: %v", msg)
	}

	type errResponse struct {
		Code    int    `json:"errorCode"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	payload := errResponse{
		Code:    code,
		Message: msg,
		Error:   "whats up doc",
	}

	respondWithJSON(w, code, payload)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Marshal JSON failed  %v", err)
		w.WriteHeader(500)

		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
