package main

import (
	"net/http"

	"github.com/Goralive/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirps(response http.ResponseWriter, request *http.Request) {
	pathVariable := request.PathValue("chirpID")
	chirpUuid, err := uuid.Parse(pathVariable)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid chirp uuid", err)
		return
	}
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.signature)
	if err != nil {
		respondWithError(response, http.StatusUnauthorized, "User not found", err)
		return
	}

	user, err := cfg.db.GetUser(request.Context(), id)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "User Not Found", err)
	}

	chirp, err := cfg.db.GetChirp(request.Context(), chirpUuid)
	if err != nil {
		respondWithError(response, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if user.ID != chirp.UserID {
		respondWithError(response, http.StatusForbidden, "Can't delete", err)
		return
	}

	deleteChirpError := cfg.db.DeleteChirp(request.Context(), chirp.ID)
	if deleteChirpError != nil {
		respondWithError(response, http.StatusInternalServerError, "Can't delete chirp", err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}
