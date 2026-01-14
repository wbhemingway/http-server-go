package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wbhemingway/http-server-go/internal/auth"
	"github.com/wbhemingway/http-server-go/internal/database"
)

const hour int = 60 * 60 * 60

func (cfg *apiConfig) loginHander(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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
	if err != nil || !passOk {
		log.Println("Error with checkpasswordhash", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("Error creating refresh token:", err)
		respondWithError(w, http.StatusInternalServerError, "Error making refresh token")
		return
	}

	dbParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	token, err := cfg.db.CreateRefreshToken(r.Context(), dbParams)
	if err != nil {
		log.Println("Error creating refresh token:", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}

	user := databaseUsertoUser(dbUser)
	respondWithJson(w, http.StatusOK, response{
		User:         user,
		Token:        accessToken,
		RefreshToken: token.Token,
	})
}
