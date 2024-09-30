package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/niccolot/Chirpy/internal/customErrors"
)

type errResponse struct {
	Error string `json:"error"`
	StatusCode int `json:"status code"`
}

type respSuccUserPostData struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

type respSuccLoginPostData struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string `json:"email"`
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
	dat, errMashal := json.Marshal(errResp)
	if errMashal != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Fatalf(fmt.Sprintf("Error marshalling JSON: %v", errMashal))
	}

	(*w).WriteHeader(err.StatusCode)
	(*w).Write(dat)
}

func respSuccesfullChirpPost(w *http.ResponseWriter, chirp *Chirp) {
	dat, errMarshal := json.Marshal(chirp)
	if errMarshal != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Fatalf(fmt.Sprintf("Error marshalling JSON: %v", errMarshal))
	}

	(*w).WriteHeader(http.StatusCreated)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullChirpsAllGet(w *http.ResponseWriter, chirps []Chirp) {
	dat, errMarshal := json.Marshal(chirps)
	if errMarshal != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Fatalf(fmt.Sprintf("Error marshalling JSON: %v", errMarshal))
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullChirpsGet(w *http.ResponseWriter, chirp *Chirp) {
	dat, errMarshal := json.Marshal(chirp)
	if errMarshal != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Fatalf(fmt.Sprintf("Error marshalling JSON: %v", errMarshal))
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullUserPost(w *http.ResponseWriter, user *User) {
	respStruct := respSuccUserPostData{
		Id: user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	dat, errMarshal := json.Marshal(respStruct)
	if errMarshal != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Fatalf(fmt.Sprintf("Error marshalling JSON: %v", errMarshal))
	}

	(*w).WriteHeader(http.StatusCreated)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullLoginPost(w *http.ResponseWriter, user *User) {
	respStruct := respSuccLoginPostData{
		Id: user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	dat, errMarshal := json.Marshal(respStruct)
	if errMarshal != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Fatalf(fmt.Sprintf("Error marshalling JSON: %v", errMarshal))
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}