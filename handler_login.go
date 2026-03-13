package main

import (
	"net/http"
	"errors"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/Karina-Pogorzelec/Chirpy/internal/auth"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email	string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
        return
    }

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		}
		return
	}
	
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't check password")
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	var ExpiresIn int
	
	if params.ExpiresInSeconds > 3600 || params.ExpiresInSeconds == 0 {
		ExpiresIn = 3600
	} else {
		ExpiresIn = params.ExpiresInSeconds
	}

	duration := time.Duration(ExpiresIn) * time.Second

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, duration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	type response struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
    Token     string    `json:"token"`
}

	respondWithJSON(w, http.StatusOK, response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}