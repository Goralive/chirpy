package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Goralive/chirpy/internal/auth"
	"github.com/Goralive/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(response http.ResponseWriter, request *http.Request) {
	type userResponse struct {
		User
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	email := params.Email
	password := params.Password

	hash_password, err := auth.HashPassword(password)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't hashed the password", err)
	}

	user, err := cfg.db.CreateUser(request.Context(), database.CreateUserParams{
		Email:          email,
		HashedPassword: hash_password,
	})
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
