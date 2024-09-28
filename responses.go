package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/niccolot/Chirpy/internal/customErrors"
)


func respondWithError(w *http.ResponseWriter, err *customErrors.CodedError) {
	message := err.Message
	code := err.StatusCode
	
	errResp := errResponse{
		Error: message,
		StatusCode: code,
	}

	fmt.Printf("error occurred: %s, status code: %d\n", message, code)
	(*w).WriteHeader(code)
	dat, e := json.Marshal(errResp)
	if e != nil {
		fmt.Printf("Error marshalling JSON: %s", e)
		return
	}

	(*w).Write(dat)
}

func respSuccesfullChirpValidation(w *http.ResponseWriter, body *string) {
	succResp := succesfullChirpPostResponse{
		CleanedBody: *body,
	}


	dat, err := json.Marshal(succResp)
	if err != nil {
		(*w).WriteHeader(500)
		fmt.Printf("Error marshalling JSON: %s", err)
		return
	}

	(*w).WriteHeader(200)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}