package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wbhemingway/http-server-go/internal/auth"
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
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Println("DecodeParams error:", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding posted json")
		return
	}
	dbParams := database.CreateChirpParams{
		Body:   cleanBody(params.Body),
		UserID: userID,
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

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println("Error Parsing uuid:", err)
		respondWithError(w, http.StatusBadRequest, "Bad request sent")
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		log.Println("Error getting uuid from db:", err)
		respondWithError(w, http.StatusNotFound, "Id not found")
		return
	}
	chirp := Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	respondWithJson(w, http.StatusOK, chirp)
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

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println("Error Parsing uuid:", err)
		respondWithError(w, http.StatusBadRequest, "Bad request sent")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirp with provided id")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirp with provided id")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
