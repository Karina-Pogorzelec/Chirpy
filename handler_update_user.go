package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Karina-Pogorzelec/Chirpy/internal/auth"
	"github.com/Karina-Pogorzelec/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
		Email	string `json:"email"`
		Password string `json:"password"`
	}

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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashedPass,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't update user")
		return
	}

	type response struct {
		ID string `json:"id"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

	respondWithJSON(w, http.StatusOK, response{
		ID: updatedUser.ID.String(),
		Email: updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}