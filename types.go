package  main

import (
    "net/http"
	"sync"
)


type apiConfig struct {
    FileserverHits int
	mu *sync.Mutex
	JwtSecret string
	PolkaApiKey string
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
	AuthorId int `json:"author_id"`
}

type succesfullUserPostResponse struct {
	Id int `json:"id"`
	Email string `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type succesfullUserPutResponse struct {
	Id int `json:"id"`
	Email string `json:"email"`
}

type succesfullLoginPostResponse struct {
	Id int `json:"id"`
	Email string `json:"email"`
	JWT string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyred bool `json:"is_chirpy_red"`
}

type succesfullRefreshPost struct {
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type polkaWebhooksPostRequest struct {
	Event string `json:"event"`
	Data polkaWebhooksData `json:"data"`
}

type polkaWebhooksData struct {
	UserId int `json:"user_id"`
}