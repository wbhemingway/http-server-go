package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/wbhemingway/http-server-go/internal/auth"
)

func (cfg *apiConfig) polkaHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token")
		return
	}

	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "You are not authorized for this request")
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Println("DecodeParams error:", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding posted json")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		log.Println("Couldn't find user:", err)
		respondWithError(w, http.StatusNotFound, "Could not find user")
	}

	w.WriteHeader(http.StatusNoContent)
}
