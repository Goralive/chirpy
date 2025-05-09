package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Goralive/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUser(request.Context(), params.UserId)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "User not found", nil)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error(), err)
	}

	chirp, err := cfg.db.CreateChirp(request.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Error during chirp save", err)
		return
	}

	respondWithJSON(response, http.StatusCreated, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func validateChirp(content string) (string, error) {
	const maxChirpLenght = 140

	if len(content) > maxChirpLenght {
		return "", errors.New("Chirp is too long")
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleaned := cleanChirp(content, profaneWords)
	return cleaned, nil
}

func cleanChirp(content string, badWords []string) string {
	words := strings.Split(content, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if slices.Contains(badWords, lower) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
