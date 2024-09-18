package main

import (
	"crypto/rand"
	"encoding/hex"
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

func postChirpHandlerWrapped(db *database.DB, cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
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
		
		dbStruct, err := db.LoadDB()
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		tokenObjPtr, errParseToken := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		if errParseToken != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("invalid token: %w", errParseToken).Error(),
				StatusCode: 403,
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

		len, err := db.GetNumChirps()
		if err != nil {
			respondWithError(&w, err)
			return 
		}

		chirp, err := db.CreateChirp(req.Body, userId)
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
		
		respSuccesfullChirpPost(&w, req.Body, id, userId)
	}
	
	return postChirpHandler
}

func getChirpsHandlerWrapped(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	getChirpsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		authorIdString := r.URL.Query().Get("author_id")
		sorting := r.URL.Query().Get("sort")
		var chirps []database.Chirp
		var err *errors.CodedError
		if authorIdString != "" {
			authorId, errAtoi := strconv.Atoi(authorIdString)
			if errAtoi != nil {
				e := errors.CodedError{
					Message: fmt.Errorf("failed to convert string to int: %w, function: %s", errAtoi, errors.GetFunctionName()).Error(),
					StatusCode: 500,
				}
				respondWithError(&w, &e)
				return
			}
			chirps, err = db.GetChirpsFromAuthor(authorId, sorting)
			if err != nil {
				respondWithError(&w, err)
				return 
			}
		} else {
			chirps, err = db.GetChirps(sorting)
			if err != nil {
				respondWithError(&w, err)
				return 
			}
		}
		
		dat, errMarshal := json.Marshal(chirps)
		if errMarshal != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to marshal json: %w, function: %s", errMarshal, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return 
		}
		
		respSuccesfullChirpGet(&w, &dat)
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

		respSuccesfullChirpGet(&w, &dat)		
	}

	return getChirpIDHandler
}

func deleteChirpIDHandlerWrapped(db *database.DB, cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	deleteChirpIDHandler := func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		tokenObjPtr, errParseToken := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		if errParseToken != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("invalid token: %w", errParseToken).Error(),
				StatusCode: 403,
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

		chirpId, errChirpId := strconv.Atoi(r.PathValue("chirpId"))
		if errChirpId != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to convert string to int: %w, function: %s", errChirpId, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}

		dbStruct, errLoad := db.LoadDB()
		if errLoad != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to load database: %w, function: %s", errGetID, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}

		if dbStruct.Users[userId].Id != dbStruct.Chirps[chirpId].AuthorId {
			e := errors.CodedError{
				Message: "invalid token, permission denied",
				StatusCode: 403,
			}
			respondWithError(&w, &e)
			return
		}

		delete(dbStruct.Chirps, chirpId)

		respSuccesfullChirpDelete(&w)
	}

	return deleteChirpIDHandler
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

		found, userIdx := dbStruct.SearchUserEmail(req.Email)
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
		if req.ExpiresInSeconds == 0 || req.ExpiresInSeconds > 60*60{
			expiresInSeconds = 60*60 // max duration of jwt is 1 hour
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

		randomSlice := make([]byte, 32)
		_, errRand := rand.Read(randomSlice)
		if errRand != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to generate refresh token: %w, function: %s", errSign, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}

		refreshToken := hex.EncodeToString(randomSlice)
		user := database.User{
			Id: userIdx,
			Email: req.Email,
			Password: req.Password,//dbStruct.Users[userIdx].Password,
			RefreshToken: refreshToken,
			IsChirpyRed: dbStruct.Users[userIdx].IsChirpyRed,

			// refresh token expires after 60 days and the date is stored as ISO 8601 format
			RefreshTokenExpiresAt: currTime.Add(60 * 24 * time.Hour).UTC().Format(time.RFC3339),
		}
		dbStruct.Users[userIdx] = user
		errWriteRefreshExp := db.WriteDB(&dbStruct)
		if errWriteRefreshExp != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to store refresh token expire time: %w, function: %s", errSign, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}

		respSuccessfullLoginPost(&w, req.Email, userIdx, signedToken, refreshToken, dbStruct.Users[userIdx].IsChirpyRed)
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
		
		errUpdate := db.UpdateUser(userId, req.Email, req.Password)
		if errUpdate != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to update database: %w, function: %s", errUpdate, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return 
		}
		respSuccesfullUserPut(&w, req.Email, userId)
	}

	return putUserHandler
}

func postRefreshHandlerWrapped(db *database.DB, cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postRefreshHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		refreshToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		dbStruct, errLoad := db.LoadDB()
		if errLoad != nil {
			respondWithError(&w, errLoad)
			return 
		}
		
		found := false
		var userIdx int
		for i, user := range(dbStruct.Users) {
			if user.RefreshToken == refreshToken {
				found = true
				userIdx = i
			}
		}

		if !found {
			e := errors.CodedError{
				Message: "refresh token does not exist",
				StatusCode: 401,
			}
			respondWithError(&w, &e)
			return
		} 

		expDate, errParse := time.Parse(time.RFC3339, dbStruct.Users[userIdx].RefreshTokenExpiresAt)
		if errParse != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to parse expiration date: %w, function: %s", errParse, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return
		}

		if time.Now().UTC().After(expDate) {
			e := errors.CodedError{
				Message: "refresh token is expired",
				StatusCode: 401,
			}
			respondWithError(&w, &e)
			return
		}
		
		currTime := time.Now().UTC()
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.RegisteredClaims{
				Issuer: "chirpy",
				IssuedAt: jwt.NewNumericDate(currTime),

				// new jwt expires after 1 hour
				ExpiresAt: jwt.NewNumericDate(currTime.Add(time.Duration(60*60)*time.Second)),
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

		respondSuccesfullRefreshPost(&w, signedToken)
	}

	return postRefreshHandler
}

func postRevokeHandlerWrapped(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	postRevokeHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		refreshToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		dbStruct, errLoad := db.LoadDB()
		if errLoad != nil {
			respondWithError(&w, errLoad)
			return 
		}
		
		found := false
		var userIdx int
		for i, user := range(dbStruct.Users) {
			if user.RefreshToken == refreshToken {
				found = true
				userIdx = i
			}
		}

		if !found {
			e := errors.CodedError{
				Message: "refresh token does not exist",
				StatusCode: 401,
			}
			respondWithError(&w, &e)
			return
		} 
		
		user := database.User{
			Id: userIdx,
			Email: dbStruct.Users[userIdx].Email,
			Password: dbStruct.Users[userIdx].Password,
			RefreshToken: "",
			RefreshTokenExpiresAt: "",
			IsChirpyRed: false,
		}

		dbStruct.Users[userIdx] = user
		errWrite := db.WriteDB(&dbStruct)
		if errWrite != nil {
			respondWithError(&w, errWrite)
			return
		}

		respSuccesfullRevokePost(&w)
	}

	return postRevokeHandler 
}

func postPolkaWebhooksHandlerWrapped(db *database.DB, cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postPolkaWebhooksHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		
		apiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
		if apiKey != cfg.PolkaApiKey {
			e := errors.CodedError{
				Message: "invalid api key",
				StatusCode: 401,
			}
			respondWithError(&w, &e)
			return 
		}
		
		decoder := json.NewDecoder(r.Body)
		req := polkaWebhooksPostRequest{}
		errDecode := decoder.Decode(&req)
		if errDecode != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to decode request: %w, function: %s", errDecode, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return 
		}

		if req.Event != "user.upgraded" {
			respSuccesfullPolkaWebhooksPost(&w)
			return
		}

		dbStruct, errLoading := db.LoadDB()
		if errLoading != nil {
			respondWithError(&w, errLoading)
			return
		}

		found, _ := dbStruct.SearchUserId(req.Data.UserId)
		if !found {
			e := errors.CodedError{
				Message: fmt.Sprintf("user_id %d not found", req.Data.UserId),
				StatusCode: 404,
			}
			respondWithError(&w, &e)
			return
		}

		errUpdate := db.UpdateSubscription(req.Data.UserId, true)
		if errUpdate != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to update database: %w, function: %s", errUpdate, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			respondWithError(&w, &e)
			return 
		}

		respSuccesfullPolkaWebhooksPost(&w)
	}

	return postPolkaWebhooksHandler
}