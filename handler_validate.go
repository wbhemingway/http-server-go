package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"slices"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type sucMsg struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, 400, "Error decoding posted json")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJson(w, 200, sucMsg{CleanedBody: cleanChirp(chirp.Body)})
}


func cleanChirp(msg string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cens := "****"
	words := strings.Split(msg, " ")
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = cens
		}
	}
	cleanedMsg := strings.Join(words, " ")
	return cleanedMsg
}