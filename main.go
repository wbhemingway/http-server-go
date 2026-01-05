package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}


func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{}
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
