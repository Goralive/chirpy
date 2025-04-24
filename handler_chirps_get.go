package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(response http.ResponseWriter, request *http.Request) {
	dbChirps, err := cfg.db.GetChirps(request.Context())
	if err != nil {
		respondWithError(response, http.StatusBadRequest, "Error during getting all the chirps", err)
		return
	}

	var chirps []Chirp

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
