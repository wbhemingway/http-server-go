package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wbhemingway/http-server-go/internal/auth"
	"github.com/wbhemingway/http-server-go/internal/database"
)

type User struct {
	Id             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
}

func (cfg *apiConfig) addUserHandler(w http.ResponseWriter, r *http.Request) {
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
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Println("Hashing password error:", err)
		respondWithError(w, http.StatusBadRequest, "Error hashing given password")
		return
	}

	dbParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	}
	dbUser, err := cfg.db.CreateUser(r.Context(), dbParams)
	if err != nil {
		log.Println("Error creating user:", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	user := User{
		Id:             dbUser.ID,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
		Email:          dbUser.Email,
		HashedPassword: dbUser.HashedPassword,
	}
	respondWithJson(w, http.StatusCreated, user)
}
