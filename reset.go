package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden access")
		return
	}
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		log.Println("Error reseting users table", err)
		respondWithError(w, 500, "error resetting users")
		return
	}

	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
}
