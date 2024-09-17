package main

import (
	"net/http"
	"sync"
    "fmt"
    "os"
    "github.com/joho/godotenv"
    "github.com/niccolot/Chirpy/internal/database"
)

func main() {
    godotenv.Load()
    jwtSecret := os.Getenv("JWT_SECRET")
    polkaApiKey := os.Getenv("POLKA_API_KEY")

    mux := http.NewServeMux()
    
    cfg := apiConfig{
        FileserverHits: 0,
        mu: &sync.Mutex{},
        JwtSecret: jwtSecret,
        PolkaApiKey: polkaApiKey,
    }

    db, err := database.NewDB("database.json")
    if err != nil {
        fmt.Println(fmt.Errorf("error creating database: %w", err).Error())
    }

    initMultiplexer(mux, &cfg, db)
    server := http.Server{
        Handler: mux,
        Addr: "localhost:8080",
    }
    server.ListenAndServe()   
}
