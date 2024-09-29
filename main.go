package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatalf(fmt.Sprintf("error loading environment variables: %v", errEnv))
	}

	dbURL := os.Getenv("DB_URL")
	db, errDB := sql.Open("postgres", dbURL)
	if errDB != nil {
		log.Fatalf(fmt.Sprintf("error creating APIconfig: %v", errDB))
	}
	
	defer db.Close()
	
	cfg, errAPIConfig := NewAPIConfig(db)
	if errAPIConfig != nil {
		log.Fatalf(fmt.Sprintf("error creating APIconfig: %v", errAPIConfig))
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr: "localhost:8080",
	}

	initMultiplexer(mux, cfg)
	server.ListenAndServe()
}