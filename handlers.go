package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/niccolot/Chirpy/internal/customErrors"
)


func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type: text/plain", "charset=utf-8")
	w.WriteHeader(200)
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
				StatusCode: 500,
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
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}
	}

	return metricsHandler
}

func resetMetricshandlerWrapperd(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	resetMetricsHandler := func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Store(0)
	}

	return resetMetricsHandler
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type: application/json", "charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		req := chirpPostRequest{}
		errDecode := decoder.Decode(&req)
		if errDecode != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to decode request: %w, function: %s", 
					errDecode, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return 
		}

		errValidation := CreateChirp(&req.Body)
		if errValidation != nil {
			respondWithError(&w, errValidation)
			return 
		}

		respSuccesfullChirpValidation(&w, &req.Body)
}