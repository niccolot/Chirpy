package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

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
		authorIdString := r.URL.Query().Get("author_id")
		sorting := r.URL.Query().Get("sort")

		var authorId uuid.UUID
		var errUUID error

		if authorIdString != "" {
			authorId, errUUID = uuid.Parse(authorIdString)
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
		}
		

		var chirpsArr []database.Chirp
		var errChirps error

		_, errSearchUser := cfg.DB.FindUserById(r.Context(), authorId)
		if errSearchUser != nil || authorIdString == "" {
			if sorting == "desc" {
				chirpsArr, errChirps = cfg.DB.GetAllChirpsDesc(r.Context())
			} else { // ASC is default option
				chirpsArr, errChirps = cfg.DB.GetAllChirpsAsc(r.Context())
			}
		} else {
			if sorting == "desc" {
				chirpsArr, errChirps = cfg.DB.GetChirpsFromAuthorDesc(r.Context(), authorId)
			} else {
				chirpsArr, errChirps = cfg.DB.GetChirpsFromAuthorAsc(r.Context(), authorId)
			}
		}

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

func deleteChirpsHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	deleteChirpsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		token, errTokenHeader := auth.GetBearerToken(r.Header)
		if errTokenHeader != nil {
			respondWithError(&w, errTokenHeader)
			return 
		}

		userId, errJWT := auth.ValidateJWT(token, cfg.JWTSecret)
		if errJWT != nil {
			respondWithError(&w, errJWT)
			return 
		}

		chirpId := r.PathValue("id")
		chirpUUID, errUUID := uuid.Parse(chirpId)
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

		chirp, errFindChirp := cfg.DB.GetChirp(r.Context(), chirpUUID)
		if errFindChirp != nil {
			e := customErrors.CodedError{
				Message: "chirp not found",
				StatusCode: http.StatusNotFound,
			}
			respondWithError(&w, &e)
			return 
		}

		errCompare := auth.CompareUUIDs(&userId, &chirp.UserID)
		if errCompare != nil {
			respondWithError(&w, errCompare)
			return 
		}	
		
		delChirpParams := &database.DeleteChirpParams{
			ID: chirp.ID,
			UserID: userId,
		}

		errDelete := cfg.DB.DeleteChirp(r.Context(), *delChirpParams)
		if errDelete != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to delete chirp %s: %w, fucntion: %s",
					string(chirpId),
					errDelete,
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}

			respondWithError(&w, &e)
			return 
		}

		respNoContent(&w)
	}

	return deleteChirpsHandler
}

func putChirpsHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	putChirpsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		token, errTokenHeader := auth.GetBearerToken(r.Header)
		if errTokenHeader != nil {
			respondWithError(&w, errTokenHeader)
			return 
		}

		userId, errJWT := auth.ValidateJWT(token, cfg.JWTSecret)
		if errJWT != nil {
			respondWithError(&w, errJWT)
			return 
		}

		decoder := json.NewDecoder(r.Body)
		req := chirpPutRequest{}
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

		chirp, errChirp := cfg.DB.GetChirp(r.Context(), req.ChirpId)
		if errChirp != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to get chirps: %w, function: %s", 
					errChirp, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		errCompare := auth.CompareUUIDs(&userId, &chirp.UserID)
		if errCompare != nil {
			respondWithError(&w, errCompare)
			return 
		}

		updateChirpParams := &database.UpdateChirpParams{
			ID: req.ChirpId,
			Body: req.Body,
		}

		errUpdate := cfg.DB.UpdateChirp(r.Context(), *updateChirpParams)
		if errUpdate != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to update chirp: %w, function: %s", 
					errUpdate, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
		}

		chirp, errFind := cfg.DB.GetChirp(r.Context(), chirp.ID)
		if errFind != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to retrieve updated chirp: %w, function: %s", 
					errFind, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		c := Chirp{}
		c.mapChirp(&chirp)

		respSuccesfullChirpPut(&w, &c)
	}

	return putChirpsHandler
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

func putUsersHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	putUsersHandlerWrapped := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		token, errTokenHeader := auth.GetBearerToken(r.Header)
		if errTokenHeader != nil {
			respondWithError(&w, errTokenHeader)
			return 
		}

		userId, errJWT := auth.ValidateJWT(token, cfg.JWTSecret)
		if errJWT != nil {
			respondWithError(&w, errJWT)
			return 
		}

		decoder := json.NewDecoder(r.Body)
		req := userPutRequest{}
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

		hashedPassword, errHash := auth.HashPassword(req.Password)
		if errHash != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to hash new password: %w, function: %s", 
					errHash, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
		}

		updateUserPars := &database.UpdateUserParams{
			ID: userId,
			Email: req.Email,
			HashedPassword: hashedPassword,
		}

		errUpdate := cfg.DB.UpdateUser(r.Context(), *updateUserPars)
		if errUpdate != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to update user: %w, function: %s", 
					errUpdate, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
		}

		user, errUser := cfg.DB.FindUserByEmail(r.Context(), req.Email)
		if errUser != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to retrieve updated user: %w, function: %s", 
					errUser, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		u := User{}
		u.mapUser(&user)

		respSuccesfullUserPut(&w, &u)
	}

	return putUsersHandlerWrapped
}

func deleteUsersHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	deleteUsersHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		token, errTokenHeader := auth.GetBearerToken(r.Header)
		if errTokenHeader != nil {
			respondWithError(&w, errTokenHeader)
			return 
		}

		userId, errJWT := auth.ValidateJWT(token, cfg.JWTSecret)
		if errJWT != nil {
			respondWithError(&w, errJWT)
			return 
		}

		userIdHeader := r.PathValue("id")
		userUUIDHeader, errUUID := uuid.Parse(userIdHeader)
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

		errCompare := auth.CompareUUIDs(&userId, &userUUIDHeader)
		if errCompare != nil {
			respondWithError(&w, errCompare)
			return 
		}	

		errDelete := cfg.DB.DeleteUser(r.Context(), userId)
		if errDelete != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to delete user: %w, fucntion: %s",
					errDelete,
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}

			respondWithError(&w, &e)
			return 
		}

		respNoContent(&w)
	}

	return deleteUsersHandler
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

		token, refreshToken, errToken := auth.MakeJWT(user.ID, cfg.JWTSecret)
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

		expiresAt := time.Now().Add(60 * 24 * time.Hour)

		refreshTokensPars := &database.CreateRefreshTokenParams{
			Token: refreshToken,
			UserID: user.ID,
			ExpiresAt: expiresAt.Format("2006-01-02 15:04:05"),
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

		respSuccesfullLoginPost(&w, &u, &token, &refreshToken)
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

		tokenObj, errObj := cfg.DB.GetRefreshToken(r.Context(), token)
		if errObj != nil {
			e := customErrors.CodedError{
				Message: "failed to retrieve refresh token from database",
				StatusCode: http.StatusNotFound,
			}
			respondWithError(&w, &e)
			return 
		}

		errValid := auth.CheckValidityRefreshToken(&tokenObj)
		if errValid != nil {
			respondWithError(&w, errValid)
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

		newToken, refreshToken, errToken := auth.MakeJWT(userId, cfg.JWTSecret)
		if errToken != nil {
			respondWithError(&w, errToken)
		}

		respSuccesfullRefreshPost(&w, newToken, refreshToken)
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

		respNoContent(&w)
	}

	return postRevokeHandler
}

func postPolkaWebhookHandlerWrapped(cfg *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	postPolkaWebhookHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type: application/json", "charset=utf-8")
		headerKey, errKey := auth.GetAPIKey(r.Header)
		if errKey != nil {
			respondWithError(&w, errKey)
			return
		}

		errCheck := auth.CheckApiKey(&headerKey, &cfg.PolkaKey)
		if errCheck != nil {
			respondWithError(&w, errCheck)
			return
		}

		decoder := json.NewDecoder(r.Body)
		req := polkaWebhookPostRequest{}
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

		if req.Event != "user.upgraded" {
			respNoContent(&w)
			return
		}

		userId := &req.Data.UserId
		_, errSearchUser := cfg.DB.FindUserById(r.Context(), *userId)
		if errSearchUser != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("user not found: %w, function: %s",
					errSearchUser, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusNotFound,
			}
			respondWithError(&w, &e)
			return 
		}

		errUpgrade := cfg.DB.UpgradeChirpyRed(r.Context(), *userId)
		if errUpgrade != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to upgrade user to chirpy red, error: %w, function: %s",
					errDecode, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}
			respondWithError(&w, &e)
			return 
		}

		respNoContent(&w)
	}

	return postPolkaWebhookHandler
}