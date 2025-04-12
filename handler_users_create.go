package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type userResponse struct {
		User
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	email := params.Email

	user, err := cfg.db.CreateUser(request.Context(), email)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't save user to db", err)
		return
	}

	log.Printf("User with email %s was saved under id: %s", user.Email, user.ID)
	respondWithJSON(response, http.StatusCreated, userResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
