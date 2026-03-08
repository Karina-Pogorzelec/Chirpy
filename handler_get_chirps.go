package main

import (
	"net/http"
)


func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	chirpsResponse := make([]Chirp, len(chirps))

	for i, dbChirp := range chirps {
		chirpsResponse[i] = Chirp{
			ID:  dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserID: dbChirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirpsResponse)
}