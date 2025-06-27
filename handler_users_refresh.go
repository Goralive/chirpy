package main

import (
	"net/http"

	"github.com/Goralive/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(response http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized)
		return
	}

	user, err = cfg.db.
}
