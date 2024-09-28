package main

import (
	"net/http"
)

func main() {
	cfg := NewAPIConfig()
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr: "localhost:8080",
	}

	initMultiplexer(mux, cfg)
	server.ListenAndServe()
}