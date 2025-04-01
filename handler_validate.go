package main

import (
	"encoding/json"
	"net/http"
)

func handlerChirpsValidate(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(500)
		response.Write([]byte("Something went wrong"))
		return
	}

	if len(params.Body) > 140 {
		resp := errorResponse{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(resp)
		if err != nil {
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(500)
			response.Write([]byte("Something went wrong"))
			return
		}
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(dat)
		return
	}

	respBody := validResponse{
		Valid: true,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(500)
		response.Write([]byte("Something went wrong"))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(dat)
}
