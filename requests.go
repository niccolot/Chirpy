package main

import "github.com/google/uuid"

type chirpPostRequest struct {
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type userPostRequest struct {
	Email string `json:"email"`
}