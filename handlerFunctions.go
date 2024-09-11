package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/niccolot/Chirpy/internal/database"
	"github.com/niccolot/Chirpy/internal/errors"
	"golang.org/x/crypto/bcrypt"
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
		req := chirpPostRequest{}
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

		len, err := db.GetNumChirps()
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
		
		respSuccesfullChirpPost(&w, req.Body, id)
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

func postUserHandlerWrapped(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	postUserHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		req := userPostRequest{}
		errDecode := decoder.Decode(&req)
		if errDecode != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to decode request: %w, function: %s", errDecode, errors.GetFunctionName()).Error(),
			}
			respondWithError(&w, &e)
			return 
		}

		user, err := db.CreateUser(req.Email, req.Password)
		if err != nil {
			respondWithError(&w, err)
			return 
		}
		
		dbStruct, err := db.LoadDB()
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		len, err := db.GetNumUsers()
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		id := len+1
		dbStruct.Users[id] = user
		errWrite := db.WriteDB(&dbStruct)
		if errWrite != nil {
			respondWithError(&w, errWrite)
			return 
		}
		
		respSuccesfullUserPost(&w, req.Email, id)
	}

	return postUserHandler
}

func postLoginHandlerWrapped(db *database.DB, cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postLoginHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		req := loginPostRequest{}
		errDecode := decoder.Decode(&req)
		if errDecode != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to decode request: %w, function: %s", errDecode, errors.GetFunctionName()).Error(),
			}
			respondWithError(&w, &e)
			return 
		}

		dbStruct, err := db.LoadDB()
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		found, userIdx, errSearchPtr := db.SearchUserEmail(req.Email)
		if errSearchPtr != nil {
			respondWithError(&w, errSearchPtr)
			return 
		}

		if !found {
			e := errors.CodedError{
				Message: fmt.Sprintf("user '%s' not found", req.Email),
				StatusCode: 404,
			}
			respondWithError(&w, &e)
			return
		}
		
		errCompPass := bcrypt.CompareHashAndPassword([]byte(dbStruct.Users[userIdx].Password), []byte(req.Password))
		if errCompPass != nil {
			e := errors.CodedError{
				Message: "unauthorized access",
				StatusCode: 401,
			}
			respondWithError(&w, &e)
			return 
		}

		var expiresInSeconds int
		if req.ExpiresInSeconds == 0 || req.ExpiresInSeconds > 24*60*60{
			expiresInSeconds = 24*60*60
		} else {
			expiresInSeconds = req.ExpiresInSeconds
		}
		
		currTime := time.Now().UTC()
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.RegisteredClaims{
				Issuer: "chirpy",
				IssuedAt: jwt.NewNumericDate(currTime),
				ExpiresAt: jwt.NewNumericDate(currTime.Add(time.Duration(expiresInSeconds)*time.Second)),
				Subject: strconv.Itoa(userIdx),
			},
		)

		signedToken, errSign := token.SignedString([]byte(cfg.JwtSecret))
		if errSign != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to sign jwt: %w, function: %s", errSign, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}

		respSuccessfullLoginPost(&w, req.Email, userIdx, signedToken)
	}

	return postLoginHandler
}

func putUserhandlerWrapped(db *database.DB, cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	putUserHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		req := userPutRequest{}
		errDecode := decoder.Decode(&req)
		if errDecode != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to decode request: %w, function: %s", errDecode, errors.GetFunctionName()).Error(),
			}
			respondWithError(&w, &e)
			return 
		}

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		tokenObjPtr, errParseToken := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		if errParseToken != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("invalid token: %w", errParseToken).Error(),
				StatusCode: 401,
			}
			respondWithError(&w, &e)
			return
		}

		userIdString, errGetID := tokenObjPtr.Claims.(*jwt.RegisteredClaims).GetSubject()		
		
		if errGetID != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("error getting user id: %w, function: %s", errGetID, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}
		
		userId, errConversion := strconv.Atoi(userIdString)
		if errConversion != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to convert userId from string to int: %w, function: %s", errConversion, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return 
		}
		db.UpdateUser(userId, req.Email, req.Password)
		respSuccesfullUserPut(&w, req.Email, userId)
	}

	return putUserHandler
}