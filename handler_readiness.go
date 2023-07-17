package main

import "net/http"

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome world"))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}
