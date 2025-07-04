package main

import (
	"encoding/json"
	"net/http"

	"github.com/Goralive/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpdateUser(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, "Provide auth token")
		return
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	email := params.Email
	password := params.Password
}
