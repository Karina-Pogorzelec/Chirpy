package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/Karina-Pogorzelec/Chirpy/internal/database"
)

var badWords = map[string]struct{}{
    "kerfuffle": {},
    "sharbert":  {},
    "fornax":    {},
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseChirpToChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body	string `json:"body"`
		UserID 	uuid.UUID `json:"user_id"`
	}	

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
        return
    }

	const maxChirpLength = 140
    if len(params.Body) > maxChirpLength {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }
    body_list := strings.Split(params.Body, " ")

    for i, word := range body_list {
        if _, exists := badWords[strings.ToLower(word)]; exists {
            body_list[i] = "****"
        }
    }

	cleanedBody := strings.Join(body_list, " ")

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseChirpToChirp(chirp))
}