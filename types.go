package main

import (
	"net/http"
	//"sync"
	"sync/atomic"
)


type apiConfig struct {
	FileserverHits atomic.Int32
	//mu *sync.Mutex
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})

	return handler
}

func NewAPIConfig() *apiConfig {
	cfg := &apiConfig{}
	cfg.FileserverHits.Store(0)

	return cfg
}

type errResponse struct {
	Error string `json:"error"`
	StatusCode int `json:"status code"`
}

type TemplateData struct {
	FileserverHits int32
}

type chirpPostRequest struct {
	Body string `json:"body"`
}

type succesfullChirpPostResponse struct {
	CleanedBody string `json:"cleaned_body"`
}