package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"github.com/niccolot/Chirpy/internal/database"
	"github.com/niccolot/Chirpy/internal/errors"
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
			e := errors.CodedError{
				Message: fmt.Errorf("error parsing template: %w, function: %s", err, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
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
		errDecode := decoder.Decode(&req)
		if errDecode != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to decode request: %w, function: %s", errDecode, errors.GetFunctionName()).Error(),
			}
			respondWithError(&w, &e)
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
		errWrite := db.WriteDB(&dbStruct)
		if errWrite != nil {
			respondWithError(&w, errWrite)
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

		dat, errMarshal := json.Marshal(chirps)
		if errMarshal != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to marshal json: %w, function: %s", err, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return 
		}
		
		respSuccesfullGet(&w, &dat)
	}

	return getChirpsHandler
}

func getChirpIDHandlerWrapped(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	getChirpIDHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to convert string to int: %w, function: %s", err, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}

		chirp, errGet := db.GetChirpID(id)
		if errGet != nil {
			respondWithError(&w, errGet)
			return
		}

		dat, errMarshal := json.Marshal(chirp)
		if errMarshal != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to marshal json: %w, function: %s", err, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return 
		}

		respSuccesfullGet(&w, &dat)		
	}

	return getChirpIDHandler
}