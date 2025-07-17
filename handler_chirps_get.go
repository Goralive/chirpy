package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(response http.ResponseWriter, request *http.Request) {
	authorId := request.URL.Query().Get("author_id")
	sortChirp := request.URL.Query().Get("sort")
	chirps := []Chirp{}
	if authorId == "" {

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

	} else {
		authorUuid, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(response, http.StatusInternalServerError, "Can't parse uuid", err)
			return
		}
		dbChirp, err := cfg.db.GetChirpByAuthor(request.Context(), authorUuid)
		if err != nil {
			respondWithError(response, http.StatusInternalServerError, "Can't get users chirps", err)
			return
		}
		for _, chirp := range dbChirp {
			chirps = append(chirps, Chirp{
				Id:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID,
			})
		}
	}

	if sortChirp == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
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
