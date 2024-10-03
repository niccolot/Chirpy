package main

import (
	"net/http"
)


func initMultiplexer(mux *http.ServeMux, cfg *apiConfig) {
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", metricshandlerWrapped(cfg))
	mux.HandleFunc("POST /admin/reset", resetMetricshandlerWrapperd(cfg))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/users", postUsersHandlerWrapped(cfg))
	mux.HandleFunc("POST /api/chirps", postChirphandlerWrapped(cfg))
	mux.HandleFunc("GET /api/chirps", getAllChirpsHandlerWrapped(cfg))
	mux.HandleFunc("GET /api/chirps/{id}", getChirspHandlerWrapped(cfg))
	mux.HandleFunc("DELETE /api/chirps/{id}", deleteChirpsHandlerWrapped(cfg))
	mux.HandleFunc("POST /api/login", postLoginHandlerWrapped(cfg))
	mux.HandleFunc("POST /api/refresh", postRefreshHandlerWrapped(cfg))
	mux.HandleFunc("POST /api/revoke", postRevokeHandlerWrapped(cfg))
	mux.HandleFunc("PUT /api/users", putUsersHandlerWrapped(cfg))
	mux.HandleFunc("POST /api/polka/webhooks", postPolkaWebhookHandlerWrapped(cfg))
}