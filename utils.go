package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
)


func respondWithError(w *http.ResponseWriter, err error) {
	(*w).WriteHeader(500)
	errResp := errResponse{
		Error: fmt.Errorf("something went wrong, %w", err).Error(),
	}

	dat, err := json.Marshal(errResp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	(*w).Write(dat)
}

func respondWithCodedError(w *http.ResponseWriter, err error) {
	//(*w).WriteHeader(500)
	//errCode := err.StatusCode
	errResp := errResponse{
		Error: fmt.Errorf("something went wrong, %w", err).Error(),
	}

	dat, err := json.Marshal(errResp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	(*w).Write(dat)
}

func respondWithSuccess(w *http.ResponseWriter, r *request) {
	succResp := succesfullResponse{
		CleanedBody: r.Body,
	}
	
	dat, err := json.Marshal(succResp)
	if err != nil {
		(*w).WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	(*w).WriteHeader(200)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullPost(w *http.ResponseWriter, body string, id int) {
	succResp := succesfullResponse{
		Id: id,
		CleanedBody: body,
	}
	
	dat, err := json.Marshal(succResp)
	if err != nil {
		(*w).WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	(*w).WriteHeader(201)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullGet(w *http.ResponseWriter, dat *[]byte) {

	(*w).WriteHeader(200)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(*dat)
}