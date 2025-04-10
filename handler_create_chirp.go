package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirps(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type chirpResponse struct {
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		UserId    uuid.UUID `json:"user_id"`
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

	user, err := cfg.db.GetUser(request.Context(), params.UserId)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "User not found", nil)
	}

	// chirp, err := cfg.db
	respondWithJSON(response, http.StatusOK, chirpResponse{
		Body: cleanChirp(content),
	})
}

func cleanChirp(content string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(content, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if slices.Contains(profaneWords, lower) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
