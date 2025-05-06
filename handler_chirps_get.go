package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(response http.ResponseWriter, request *http.Request) {
	chirps := []Chirp{}
	dbChirps, err := cfg.db.GetChirps(request.Context())
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Error during getting all the chirps", err)
		return
	}

	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	respondWithJSON(response, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(response http.ResponseWriter, request *http.Request) {
	chirpStringUuid := request.PathValue("chirpID")

	chirpUuid, err := uuid.Parse(chirpStringUuid)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Invalid chirp uuid", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(request.Context(), chirpUuid)
	if err != nil {
		respondWithError(response, http.StatusNotFound, "Chirp Not found", err)
		return
	}

	respondWithJSON(response, http.StatusOK, Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	})
}
