package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/niccolot/Chirpy/internal/customErrors"
	"github.com/niccolot/Chirpy/internal/database"
)


func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type: text/plain", "charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func metricshandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	metricsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: text/html", "charset=utf-8")
		tmpl, err := template.ParseFiles("index_admin.html")
		if err != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("internal Server Error: %w, function: %s", 
					err, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return
		}
		
		data := &TemplateData{
			FileserverHits: cfg.FileserverHits.Load(),
		}

		err = tmpl.Execute(w, *data)
		if err != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("error parsing template: %w, function: %s", 
					err, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return
		}
	}

	return metricsHandler
}

func resetMetricshandlerWrapperd(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	resetMetricsHandler := func(w http.ResponseWriter, r *http.Request) {
		if cfg.Platform != "dev" {
			e := customErrors.CodedError{
				Message: "forbidden request",
				StatusCode: http.StatusForbidden,
			}
			respondWithError(&w, &e)
			return
		}

		errDelete := cfg.DB.Reset(r.Context())
		if errDelete != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("error executing reset request: %w, function: %s", 
					errDelete,
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return
		}
	}

	return resetMetricsHandler
}

func postChirphandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postChirpHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		req := chirpPostRequest{}
		errDecode := decoder.Decode(&req)
		if errDecode != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to decode request: %w, function: %s", 
					errDecode, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		errValidation := ValidateChirp(&req.Body)
		if errValidation != nil {
			respondWithError(&w, errValidation)
			return 
		}

		chirpparams := database.CreateChirpParams{
			Body: req.Body,
			UserID: req.UserId,
		}

		chirp, errChirp := cfg.DB.CreateChirp(r.Context(), chirpparams)
		if errChirp != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to create chirp: %w, function: %s", 
					errChirp, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		c := Chirp{}
		c.mapChirp(&chirp)

		respSuccesfullChirpPost(&w, &c)
	}

	return postChirpHandler
} 

func getAllChirpsHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	getAllChirpsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		chirpsArr, errChirps := cfg.DB.GetAllChirps(r.Context())
		if errChirps != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to get chirps: %w, function: %s", 
					errChirps, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		cArr := make([]Chirp, len(chirpsArr))
		for i, c := range chirpsArr {
			cArr[i].mapChirp(&c)
		}

		respSuccesfullChirpsAllGet(&w, cArr)
	}

	return getAllChirpsHandler
}

func getChirspHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	getChirpsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		id := r.PathValue("id")
		uuid, errUUID := uuid.Parse(id)
		if errUUID != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("error parsing uuid: %w, function: %s", 
					errUUID, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		chirp, errChirp := cfg.DB.GetChirp(r.Context(), uuid)
		if errChirp != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to get chirp: %w, function: %s", 
					errUUID, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusNotFound,
			}
			respondWithError(&w, &e)
			return 
		}

		c := Chirp{}
		c.mapChirp(&chirp)

		respSuccesfullChirpsGet(&w, &c)
	}

	return getChirpsHandler
}

func postUsersHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postUsersHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		req := userPostRequest{}
		errDecode := decoder.Decode(&req)
		if errDecode != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to decode request: %w, function: %s", 
					errDecode, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		user, errUser := cfg.DB.CreateUser(r.Context(), req.Email)
		if errUser != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to create user: %w, function: %s", 
					errUser, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		u := User{}
		u.mapUser(&user)

		respSuccesfullUserPost(&w, &u)
	}

	return postUsersHandler
}