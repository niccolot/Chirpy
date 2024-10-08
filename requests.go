package main

import "github.com/google/uuid"

type chirpPostRequest struct {
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type chirpPutRequest struct {
	ChirpId uuid.UUID `json:"id"`
	Body string `json:"body"`
}

type userPostRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type userPutRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type loginPostRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type polkaWebhookPostRequest struct {
	Event string `json:"event"`
	Data polkaWebhookData `json:"data"`
}

type polkaWebhookData struct {
	UserId uuid.UUID `json:"user_id"`
}