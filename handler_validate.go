package main

import (
	"net/http"
	"encoding/json"
    "strings"
)

var badWords = map[string]struct{}{
    "kerfuffle": {},
    "sharbert":  {},
    "fornax":    {},
}


func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
        Body string `json:"body"`
    }
    type returnVals struct {
    CleanedBody string `json:"cleaned_body"`
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

    respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: strings.Join(body_list, " ")})
}
