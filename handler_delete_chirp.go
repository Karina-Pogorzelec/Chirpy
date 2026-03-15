package main

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/Karina-Pogorzelec/Chirpy/internal/auth"
	"github.com/Karina-Pogorzelec/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	chirpID := r.PathValue("chirpID")
	parsedChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), parsedChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can only delete your own chirps")
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: parsedChirpID,
		UserID:  userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}