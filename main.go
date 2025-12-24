package main

import "net/http"

func main() {
	const filepathRoot = "."
	const port = "8080"
	
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	
	serve := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	
	serve.ListenAndServe()
}
