package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/niccolot/Chirpy/internal/auth"
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

		token, errGetToken := auth.GetBearerToken(r.Header)
		if errGetToken != nil {
			respondWithError(&w, errGetToken)
			return 
		}

		id, errValidateAuthor := auth.ValidateJWT(token, cfg.JWTSecret)
		if errValidateAuthor != nil {
			respondWithError(&w, errValidateAuthor)
		}

		errChirpValidation := ValidateChirp(&req.Body)
		if errChirpValidation != nil {
			respondWithError(&w, errChirpValidation)
			return 
		}

		chirpPars := database.CreateChirpParams{
			Body: req.Body,
			UserID: id,
		}

		chirp, errChirp := cfg.DB.CreateChirp(r.Context(), chirpPars)
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

		hashed_password, errHashing := auth.HashPassword(req.Password)
		if errHashing != nil {
			respondWithError(&w, errHashing)
			return
		}

		userPars := &database.CreateUserParams{
			Email: req.Email,
			HashedPassword: hashed_password,
		}

		user, errUser := cfg.DB.CreateUser(r.Context(), *userPars)
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

func postLoginHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postLoginhandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		decoder := json.NewDecoder(r.Body)
		req := loginPostRequest{}
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

		user, errUser := cfg.DB.FindUserByEmail(r.Context(), req.Email)
		if errUser != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to find user: %w, function: %s", 
					errUser, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		check := auth.CheckPasswordHash(req.Password, user.HashedPassword)
		if check != nil {
			respondWithError(&w, check)
			return 
		}

		token, errToken := auth.MakeJWT(user.ID, cfg.JWTSecret)
		if errToken != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to generate jwt: %w, function: %s", 
					errToken, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		refreshTokenString, errRefresh := auth.MakeRefreshToken()
		if errRefresh != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to generate refresh token string: %w, function: %s", 
					errRefresh, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		refreshTokensPars := &database.CreateRefreshTokenParams{
			Token: refreshTokenString,
			UserID: user.ID,
		}

		_, errRefreshObj := cfg.DB.CreateRefreshToken(r.Context(), *refreshTokensPars)
		if errRefreshObj != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to generate refresh token object: %w, function: %s", 
					errRefreshObj, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		u := User{}
		u.mapUser(&user)

		respSuccesfullLoginPost(&w, &u, &token, &refreshTokenString)
	}

	return postLoginhandler
}

func postRefreshHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postRefreshHandler := func(w http.ResponseWriter, r *http.Request) {
		token, errHeader := auth.GetBearerToken(r.Header)
		if errHeader != nil {
			respondWithError(&w, errHeader)
			return
		}

		userId, errSearch := cfg.DB.GetUserFromRefreshToken(r.Context(), token)
		if errSearch != nil {
			e := customErrors.CodedError{
				Message: "invalid jwt",
				StatusCode: http.StatusUnauthorized,
			}
			respondWithError(&w, &e)
			return 
		}

		newToken, errToken := auth.MakeJWT(userId, cfg.JWTSecret)
		if errToken != nil {
			respondWithError(&w, errToken)
		}

		respSuccesfullRefreshPost(&w, newToken)
	}

	return postRefreshHandler
}

func postRevokeHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postRevokeHandler := func(w http.ResponseWriter, r *http.Request) {
		token, errHeader := auth.GetBearerToken(r.Header)
		if errHeader != nil {
			respondWithError(&w, errHeader)
			return
		}

		errRevoke := cfg.DB.RevokeToken(r.Context(), token)
		if errRevoke != nil {
			e := customErrors.CodedError{
				Message: "token not in database",
				StatusCode: http.StatusNotFound,
			}
			respondWithError(&w, &e)
			return 
		}

		respSuccesfullRevokePost(&w)
	}

	return postRevokeHandler
}