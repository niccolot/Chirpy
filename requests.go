package main

import "github.com/google/uuid"

type chirpPostRequest struct {
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type userPostRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type loginPostRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ExpiresInSeconds int `json:"expires_in_seconds"`
}