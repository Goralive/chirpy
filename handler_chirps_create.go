package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Goralive/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
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
		return
	}

	chirp, err := cfg.db.CreateChirp(request.Context(), database.CreateChirpParams{
		Body:   cleanChirp(content),
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Error during chirp save", err)
		return
	}

	// chirp, err := cfg.db
	respondWithJSON(response, http.StatusCreated, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
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
