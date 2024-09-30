package auth

import (
	"fmt"
	"net/http"

	"github.com/niccolot/Chirpy/internal/customErrors"
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