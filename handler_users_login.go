package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Goralive/chirpy/internal/auth"
	"github.com/Goralive/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLoginUser(response http.ResponseWriter, request *http.Request) {
	errorMessage := "Incorrect email or password"

	type userResponse struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	passwordValidationError := auth.CheckPasswordHash(user.HashedPassword, params.Password)

	if passwordValidationError != nil {
		respondWithError(response, http.StatusUnauthorized, errorMessage, passwordValidationError)
		return
	}

	tokenExpiration := time.Duration(3600) * time.Second

	token, err := auth.MakeJWT(user.ID, cfg.signature, tokenExpiration)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour)
	data, err := cfg.db.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(response, http.StatusOK, userResponse{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token:        token,
		RefreshToken: data.Token,
	})
}
