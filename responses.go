package main

import (
	"encoding/json"
	"fmt"
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
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type respSuccUserPutData struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type respSuccLoginPostData struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type respSuccRefreshPostData struct {
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func respondWithError(w *http.ResponseWriter, err *customErrors.CodedError) {
	message := err.Message
	code := err.StatusCode
	
	errResp := errResponse{
		Error: message,
		StatusCode: code,
	}

	fmt.Printf("error occurred: %s, status code: %d\n", message, code)
	dat, errMarshal := json.Marshal(errResp)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
	}

	(*w).WriteHeader(err.StatusCode)
	(*w).Write(dat)
}

func respSuccesfullChirpPost(w *http.ResponseWriter, chirp *Chirp) {
	dat, errMarshal := json.Marshal(chirp)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
	}

	(*w).WriteHeader(http.StatusCreated)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullChirpPut(w *http.ResponseWriter, chirp *Chirp) {
	dat, errMarshal := json.Marshal(chirp)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullChirpsAllGet(w *http.ResponseWriter, chirps []Chirp) {
	dat, errMarshal := json.Marshal(chirps)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullChirpsGet(w *http.ResponseWriter, chirp *Chirp) {
	dat, errMarshal := json.Marshal(chirp)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
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
		IsChirpyRed: user.IsChirpyred,
	}

	dat, errMarshal := json.Marshal(respStruct)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
	}

	(*w).WriteHeader(http.StatusCreated)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullUserPut(w *http.ResponseWriter, user *User) {
	respStruct := respSuccUserPutData{
		Id: user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyred,
	}

	dat, errMarshal := json.Marshal(respStruct)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullLoginPost(w *http.ResponseWriter, user *User, jwt *string, refreshToken *string) {
	respStruct := respSuccLoginPostData{
		Id: user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: *jwt,
		RefreshToken: *refreshToken,
		IsChirpyRed: user.IsChirpyred,
	}

	dat, errMarshal := json.Marshal(respStruct)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respSuccesfullRefreshPost(w *http.ResponseWriter, token string, refreshToken string) {
	respStruct := respSuccRefreshPostData{
		Token: token,
		RefreshToken: refreshToken,
	}

	dat, errMarshal := json.Marshal(respStruct)
	if errMarshal != nil {
		customErrors.ErrorMarshal(w, errMarshal)
		return 
	}

	(*w).WriteHeader(http.StatusOK)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(dat)
}

func respNoContent(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusNoContent)
}