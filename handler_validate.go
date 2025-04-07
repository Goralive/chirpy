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
	const maxChirpLenght = 140

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > maxChirpLenght {
		respondWithError(response, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(response, http.StatusOK, validResponse{
		Valid: true,
	})
}
