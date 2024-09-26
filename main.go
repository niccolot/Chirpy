package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/niccolot/Chirpy/internal/database"
	"github.com/niccolot/Chirpy/internal/jsondatabase"
)

func main() {
    godotenv.Load()
    jwtSecret := os.Getenv("JWT_SECRET")
    polkaApiKey := os.Getenv("POLKA_API_KEY")
    dbURL := os.Getenv("DB_URL")

    db, errDB := sql.Open("postgres", dbURL)
    if errDB != nil {
        log.Fatalf("error creating database: %v", errDB)
    }

    mux := http.NewServeMux()
    
    cfg := apiConfig{
        FileserverHits: 0,
        mu: &sync.Mutex{},
        JwtSecret: jwtSecret,
        PolkaApiKey: polkaApiKey,
        dbQueries: database.New(db),
    }

    debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

    if *debug {
        fmt.Println("Debug mode enabled: Deleting database.json...")
		err := os.Remove("database.json")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("database.json does not exist, nothing to delete.")
			} else {
				log.Fatalf("Error deleting database.json: %v", err)
			}
		} else {
			fmt.Println("database.json successfully deleted.")
		}
    }

    jsondb, errJSONDB := jsondatabase.NewDB("database.json")
    if errJSONDB != nil {
        log.Fatalf("error creating database: %v", errJSONDB)
    }

    initMultiplexer(mux, &cfg, jsondb)
    server := http.Server{
        Handler: mux,
        Addr: "localhost:8080",
    }
    server.ListenAndServe()   
}
