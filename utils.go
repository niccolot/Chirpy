package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/niccolot/Chirpy/internal/errors"
)


func respondWithError(w *http.ResponseWriter, err *errors.CodedError) {
	message := err.Message
	code := err.StatusCode
	
	errResp := errResponse{
		Error: message,
		StatusCode: code,
	}

	fmt.Printf("error occurred: %s, status code: %d", message, code)
	(*w).WriteHeader(code)
	dat, e := json.Marshal(errResp)
	if e != nil {
		fmt.Printf("Error marshalling JSON: %s", e)
		return
	}

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
		fmt.Printf("Error marshalling JSON: %s", err)
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