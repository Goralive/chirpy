package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

	content := params.Body

	if len(content) > maxChirpLenght {
		respondWithError(response, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(response, http.StatusOK, validResponse{
		Valid: true,
	})
}

func cleanChirp(content string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(content, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		for _, profane := range profaneWords {
			if lower == profane {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, "")
}
