package main

import (
	"net/http"
	"sync"
    "fmt"
    "github.com/niccolot/Chirpy/internal/database"
)

func main() {
    mux := http.NewServeMux()
    
    cfg := apiConfig{
        FileserverHits: 0,
        mu: &sync.Mutex{},
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
