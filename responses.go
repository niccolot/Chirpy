package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/niccolot/Chirpy/internal/customErrors"
)

type errResponse struct {
	Error string `json:"error"`
	StatusCode int `json:"status code"`
}

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
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error marshalling JSON: %s", e)
		return
	}

	(*w).WriteHeader(err.StatusCode)
	(*w).Write(dat)
}

func respSuccesfullChirpValidation(w *http.ResponseWriter, body *string) {
	succResp := succesfullChirpPostResponse{
		CleanedBody: *body,
	}

	dat, errMarshal := json.Marshal(succResp)
	if errMarshal != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error marshalling JSON: %s", errMarshal)
		return
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullUserPost(w *http.ResponseWriter, user *User) {
	dat, errMarshal := json.Marshal(user)
	if errMarshal != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error marshalling JSON: %s", errMarshal)
		return
	}

	(*w).WriteHeader(http.StatusCreated)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)

}