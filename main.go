package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	
	serve := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	serve.ListenAndServe()
}