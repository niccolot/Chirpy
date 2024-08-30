package main

import (
	"net/http"
	"html/template"
    "log"
	"encoding/json"
	"github.com/niccolot/Chirpy/internal/database"
)


func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type: text/plain", "charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func metricsHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {

	metricsHandler := func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: text/html", "charset=utf-8")
	
		tmpl, err := template.ParseFiles("index_admin.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error parsing template:", err)
			return
		}
	
		err = tmpl.Execute(w, cfg)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error executing template:", err)
			return
		}
	}

	return metricsHandler
}

func postChirpHandlerWrapped(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	postChirpHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		req := request{}
		err := decoder.Decode(&req)
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		chirp, err := db.CreateChirp(req.Body)
		if err != nil {
			respondWithError(&w, err)
			return 
		}
		
		dbStruct, err := db.LoadDB()
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		len, err := db.GetDBLength()
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		id := len+1
		dbStruct.Chirps[id] = chirp
		err = db.WriteDB(&dbStruct)
		if err != nil {
			respondWithError(&w, err)
			return 
		}
		
		respSuccesfullPost(&w, req.Body, id)
	}
	
	return postChirpHandler
}

func getChirpsHandlerWrapped(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	getChirpsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")

		chirps, err := db.GetChirps()
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		dat, err := json.Marshal(chirps)
		if err != nil {
			respondWithError(&w, err)
			return 
		}
		
		respSuccesfullGet(&w, &dat)
	}

	return getChirpsHandler
}