package main

import (
	"net/http"
	"time"

	"github.com/Goralive/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(response http.ResponseWriter, request *http.Request) {
	type tokenResponse struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, "No token in header", err)
		return
	}

	dbRefreshToken, err := cfg.db.GetUserFromRefreshToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, "User not found from refreshToken", err)
		return
	}

	tokenExpiration := time.Duration(3600) * time.Second
	token, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.signature, tokenExpiration)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Can't create jwt token", err)
		return
	}

	respondWithJSON(response, http.StatusOK, tokenResponse{
		Token: token,
	})
}

func (cfg *apiConfig) handlerRevokeToken(response http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, "No token in header", err)
		return
	}

	dbRefreshToken, err := cfg.db.GetUserFromRefreshToken(request.Context(), refreshToken)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, "User not found from refreshToken", err)
		return
	}

	errors := cfg.db.RevokeToken(request.Context(), dbRefreshToken.Token)
	if errors != nil {
		respondWithError(response, http.StatusInternalServerError, "Can't revoke token", err)
	}

	response.WriteHeader(http.StatusNoContent)
}
