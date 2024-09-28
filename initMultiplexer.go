package main

import (
	"net/http"
)


func initMultiplexer(mux *http.ServeMux, cfg *apiConfig) {
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", metricshandlerWrapped(cfg))
	mux.HandleFunc("POST /admin/reset", resetMetricshandlerWrapperd(cfg))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
}