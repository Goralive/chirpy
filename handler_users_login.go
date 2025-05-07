package main

import (
	"encoding/json"
	"net/http"

	"github.com/Goralive/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLoginUser(response http.ResponseWriter, request *http.Request) {
	errorMessage := "Incorrect email or password"

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
		respondWithError(response, http.StatusInternalServerError, "Coudn't decode params", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(request.Context(), params.Email)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, errorMessage, err)
		return
	}

	isCorrectPassword := auth.CheckPasswordHash(user.HashedPassword, params.Password)

	if isCorrectPassword != nil {
		respondWithError(response, http.StatusUnauthorized, errorMessage, err)
		return
	}

	respondWithJSON(response, http.StatusOK, userResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
