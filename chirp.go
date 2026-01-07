package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wbhemingway/http-server-go/internal/database"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) addChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("DecodeParams error:", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding posted json")
		return
	}
	dbParams := database.CreateChirpParams{
		Body:   cleanBody(params.Body),
		UserID: params.UserId,
	}
	dbChirp, err := cfg.db.CreateChirp(r.Context(), dbParams)
	if err != nil {
		// add proper error checking to see if it's 400 or 500 error
		log.Println("Error creating chirp:", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	chirp := Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	respondWithJson(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Println("Error getting users:", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting users")
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, chirp := range dbChirps {
		chirps[i] = Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
	}
	respondWithJson(w, http.StatusOK, chirps)
}

func cleanBody(msg string) string {
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
