package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wbhemingway/http-server-go/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error making database conection")
	}

	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"

	apiCfg := &apiConfig{
		queries: dbQueries,
	}
	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		http.StripPrefix(
			"/app",
			apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("."))),
		),
	)
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	serve := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Chirpy Server Started!")
	serve.ListenAndServe()
	fmt.Println("Chirpy Server Stopped!")
}
