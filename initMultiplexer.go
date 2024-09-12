package main

import (
	"net/http"
    "github.com/niccolot/Chirpy/internal/database"
)


func initMultiplexer(mux *http.ServeMux, cfg *apiConfig, db *database.DB) {
	mux.Handle("/app/*", http.StripPrefix("/app/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
    mux.HandleFunc("GET /api/healthz", healthzHandler)
    mux.HandleFunc("GET /admin/metrics/*", metricsHandlerWrapped(cfg))
    mux.HandleFunc("/api/reset", func (w http.ResponseWriter, r *http.Request) {
        cfg.FileserverHits = 0
    })
    mux.HandleFunc("POST /api/chirps", postChirpHandlerWrapped(db))
    mux.HandleFunc("GET /api/chirps", getChirpsHandlerWrapped(db))
    mux.HandleFunc("GET /api/chirps/{id}", getChirpIDHandlerWrapped(db))
    mux.HandleFunc("POST /api/users", postUserHandlerWrapped(db))
    mux.HandleFunc("PUT /api/users", putUserhandlerWrapped(db, cfg))
    mux.HandleFunc("POST /api/login", postLoginHandlerWrapped(db, cfg))
    mux.HandleFunc("POST /api/refresh", postRefreshHandlerWrapped(db, cfg))
    mux.HandleFunc("POST /api/revoke", postRevokeHandlerWrapped(db))
} 