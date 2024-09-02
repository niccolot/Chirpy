package  main

import (
    "net/http"
	"sync"
)

type apiConfig struct {
    FileserverHits int
	mu *sync.Mutex
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mu.Lock()
		cfg.FileserverHits += 1
		cfg.mu.Unlock()
		next.ServeHTTP(w,r)
	})

	return handler
}

type request struct {
	Body string `json:"body"`
}

type errResponse struct {
	Error string `json:"error"`
	StatusCode int `json:"status code"`
}

type succesfullResponse struct {
	Id int `json:"id"`
	CleanedBody string `json:"body"`
}