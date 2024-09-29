package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/niccolot/Chirpy/internal/customErrors"
	"github.com/niccolot/Chirpy/internal/database"
)


type apiConfig struct {
	DB *database.Queries
	FileserverHits atomic.Int32
	Platform string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})

	return handler
}

func NewAPIConfig(db *sql.DB) (*apiConfig, *customErrors.CodedError) {
	errEnv := godotenv.Load()
	if errEnv != nil {
		e := customErrors.CodedError{
			Message: fmt.Errorf("error loading environment variables: %w, function: %s",
				errEnv,
				customErrors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return &apiConfig{}, &e
	}
	
	cfg := &apiConfig{}
	cfg.FileserverHits.Store(0)
	dbQueries := database.New(db)
	cfg.DB = dbQueries
	platform := os.Getenv("PLATFORM")
	cfg.Platform = platform

	return cfg, nil
}