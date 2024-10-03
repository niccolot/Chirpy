package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/niccolot/Chirpy/internal/customErrors"
	"github.com/niccolot/Chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, *customErrors.CodedError) {
	password_bytes := []byte(password)
	hash, errHashing := bcrypt.GenerateFromPassword(password_bytes, bcrypt.DefaultCost)
	if errHashing != nil {
		e := customErrors.CodedError{
			Message: fmt.Errorf("error hashing password: %w, function: %s", errHashing, customErrors.GetFunctionName()).Error(),
			StatusCode: http.StatusInternalServerError,
		}
		return "", &e
	}

	return string(hash), nil
} 

func CheckPasswordHash(password string, hash string) *customErrors.CodedError {
	errCompPass := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errCompPass != nil {
		e := customErrors.CodedError{
			Message: "Incorrect email or password",
			StatusCode: http.StatusUnauthorized,
		}
		return &e
	}

	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, *customErrors.CodedError) {
	currTime := time.Now().UTC()
	expiresIn := 60*60 // 1 hour jwt duration
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer: "chirpy",
			IssuedAt: jwt.NewNumericDate(currTime),
			ExpiresAt: jwt.NewNumericDate(currTime.Add(time.Duration(expiresIn)*time.Second)),
			Subject: string(userID.String()),
		},
	)

	signedToken, errSign := token.SignedString([]byte(tokenSecret))
	if errSign != nil {
		e := customErrors.CodedError{
			Message: fmt.Errorf("failed to sign jwt: %w, function: %s", 
				errSign, 
				customErrors.GetFunctionName()).Error(),
			StatusCode: http.StatusInternalServerError,
		}
		
		return "", &e
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, *customErrors.CodedError) {
	token, errParseToken := jwt.ParseWithClaims(tokenString, 
		&jwt.RegisteredClaims{}, 
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				errMethod := customErrors.CodedError{
					Message: fmt.Errorf("unexpected signing method: %v, function: %s", 
						token.Header["alg"],
						customErrors.GetFunctionName()).Error(),
					StatusCode: http.StatusInternalServerError,
				}

				return uuid.UUID{}, &errMethod
			}

			return []byte(tokenSecret), nil
		})

	if errParseToken != nil {
		e := customErrors.CodedError{
			Message: fmt.Errorf("failed to parse token: %w, function: %s", 
				errParseToken,
				customErrors.GetFunctionName()).Error(),
			StatusCode: http.StatusUnauthorized,
		}
		return uuid.UUID{}, &e
	}

	errValid := token.Claims.Valid()
	if errValid != nil {
		e := customErrors.CodedError{
			Message: fmt.Errorf("invalid token: %w, function: %s", 
				errValid,
				customErrors.GetFunctionName()).Error(),
			StatusCode: http.StatusUnauthorized,
		}

		return uuid.UUID{}, &e
	}

	id, errParseUUID := uuid.Parse(token.Claims.(*jwt.RegisteredClaims).Subject)
    if errParseUUID != nil {
        e := customErrors.CodedError{
			Message: fmt.Errorf("failed to parse string into UUID: %w, function: %s", 
				errValid,
				customErrors.GetFunctionName()).Error(),
			StatusCode: http.StatusInternalServerError,
		}

		return uuid.UUID{}, &e
    }

	 return  id, nil	
}

func GetBearerToken(headers http.Header) (string, *customErrors.CodedError) {
	token := strings.TrimPrefix(headers.Get("Authorization"), "Bearer ")
	if token == "" {
		e := customErrors.CodedError{
			Message: "request header must contain the jwt",
			StatusCode: http.StatusUnauthorized,
		}

		return "", &e
	}

	return token, nil
}

func MakeRefreshToken() (string, *customErrors.CodedError) {
	randomSlice := make([]byte, 32)
		_, errRand := rand.Read(randomSlice)
		if errRand != nil {
			e := customErrors.CodedError{
				Message: fmt.Errorf("failed to generate refresh token: %w, function: %s", 
					errRand, 
					customErrors.GetFunctionName()).Error(),
				StatusCode: http.StatusInternalServerError,
			}

			return "", &e
		}

		refreshToken := hex.EncodeToString(randomSlice)

		return refreshToken, nil
}

func CheckValidityRefreshToken(tokenObj *database.RefreshToken) *customErrors.CodedError {
	notExpired := time.Now().Format("2006-01-02 15:04:05") <= tokenObj.ExpiresAt
	notRevoked :=  !tokenObj.RevokedAt.Valid

	valid := notExpired && notRevoked

	if !valid {
		e := &customErrors.CodedError{
			Message: "invalid refresh token",
			StatusCode: http.StatusUnauthorized,
		}

		return e
	}

	return nil
} 

func CompareUUIDs(uuid1 *uuid.UUID, uuid2 *uuid.UUID) *customErrors.CodedError {
	if *uuid1 != *uuid2 {
		e := &customErrors.CodedError{
			Message: "invalid user",
			StatusCode: http.StatusForbidden,
		}

		return e
	}

	return nil
}

func GetAPIKey(headers http.Header) (string, *customErrors.CodedError) {
	key := strings.TrimPrefix(headers.Get("Authorization"), "ApiKey ")
	if key == "" {
		e := customErrors.CodedError{
			Message: "request header must contain the polka api key",
			StatusCode: http.StatusUnauthorized,
		}

		return "", &e
	}

	return key, nil
}

func CheckApiKey(key1 *string, key2 *string) *customErrors.CodedError {
	if *key1 != *key2 {
		e := &customErrors.CodedError{
			Message: "invalid key",
			StatusCode: http.StatusUnauthorized,
		}

		return e
	}

	return nil
}