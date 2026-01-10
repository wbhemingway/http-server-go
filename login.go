package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wbhemingway/http-server-go/internal/auth"
)

func (cfg *apiConfig) loginHander(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("DecodeParams error:", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding posted json")
		return
	}
	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Println("Error finding the user:", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	
	passOk, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		log.Println("Error with checkpasswordhash", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	
	if !passOk {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	user := User{
		Id:             dbUser.ID,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
		Email:          dbUser.Email,
	}
	respondWithJson(w, http.StatusOK, user)
}