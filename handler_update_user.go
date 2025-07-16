package main

import (
	"encoding/json"
	"net/http"

	"github.com/Goralive/chirpy/internal/auth"
	"github.com/Goralive/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateUser(response http.ResponseWriter, request *http.Request) {
	type userResponse struct {
		User
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, "Provide auth token", err)
		return
	}

	id, jwtError := auth.ValidateJWT(token, cfg.signature)

	if jwtError != nil {
		respondWithError(response, http.StatusUnauthorized, "invalid token", jwtError)
		return
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	decodeErr := decoder.Decode(&params)
	if decodeErr != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't decode parameters", decodeErr)
		return
	}
	email := params.Email
	password := params.Password

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Could't hashed the password", err)
	}

	user, err := cfg.db.UpdateUser(request.Context(), database.UpdateUserParams{
		ID:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't save user to db", err)
		return
	}

	respondWithJSON(response, http.StatusOK, userResponse{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}

func (cfg *apiConfig) handlerChirpyRedWebhook(response http.ResponseWriter, request *http.Request) {
	type data struct {
		UserId uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	decoderErr := decoder.Decode(&params)
	if decoderErr != nil {
		respondWithError(response, http.StatusInternalServerError, "Can't parsed bodyo", decoderErr)
		return
	}

	event := params.Event
	userId := params.Data.UserId

	if event != "user.upgraded" {
		response.WriteHeader(http.StatusNoContent)
		return
	}

	err := cfg.db.UpdateUserChirpyRed(request.Context(), userId)
	if err != nil {
		respondWithError(response, http.StatusNotFound, "", err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}
