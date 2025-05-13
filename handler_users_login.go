package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Goralive/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLoginUser(response http.ResponseWriter, request *http.Request) {
	errorMessage := "Incorrect email or password"

	type userResponse struct {
		User
		Token string `json:"token"`
	}

	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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

	passwordValidationError := auth.CheckPasswordHash(user.HashedPassword, params.Password)

	if passwordValidationError != nil {
		respondWithError(response, http.StatusUnauthorized, errorMessage, passwordValidationError)
		return
	}

	expiration := params.ExpiresInSeconds

	if expiration == 0 || expiration > 3600 {
		expiration = 3600
	}

	tokenExpiration := time.Duration(expiration) * time.Second

	token, err := auth.MakeJWT(user.ID, cfg.signature, tokenExpiration)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't create token", err)
	}

	respondWithJSON(response, http.StatusOK, userResponse{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	})
}
