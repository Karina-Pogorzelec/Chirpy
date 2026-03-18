package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/Karina-Pogorzelec/Chirpy/internal/database"
)


func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	
	s := r.URL.Query().Get("author_id")
	var err error
	var dbChirps []database.Chirp

	if s == "" {
		dbChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
			return
		}
	} else {
		authorID, parseErr := uuid.Parse(s)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}

		dbChirps, err = cfg.db.GetChirpsForAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps for author")
			return
		}
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
        chirps = append(chirps, Chirp{
            ID:        dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            Body:      dbChirp.Body,
            UserID:    dbChirp.UserID,
        })
    }

	respondWithJSON(w, http.StatusOK, chirps)
}