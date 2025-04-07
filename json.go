package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(response http.ResponseWriter, code int, message string, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", message)
	}
	respondWithJSON(response, code, errorResponse{
		Error: message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
