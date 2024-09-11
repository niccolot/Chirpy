package  main

import (
    "net/http"
	"sync"
)


type apiConfig struct {
    FileserverHits int
	mu *sync.Mutex
	JwtSecret string
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

type chirpPostRequest struct {
	Body string `json:"body"`
}

type userPostRequest struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

type userPutRequest struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

type loginPostRequest struct {
	Password string `json:"password"`
	Email string `json:"email"`
	ExpiresInSeconds int `json:"expires_in_seconds"`
}

type errResponse struct {
	Error string `json:"error"`
	StatusCode int `json:"status code"`
}

type succesfullChirpPostResponse struct {
	Id int `json:"id"`
	CleanedBody string `json:"body"`
}

type succesfullUserPostResponse struct {
	Id int `json:"id"`
	Email string `json:"email"`
}

type succesfullUserPutResponse struct {
	Id int `json:"id"`
	Email string `json:"email"`
}

type succesfullLoginPostResponse struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}