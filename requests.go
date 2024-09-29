package main

type chirpPostRequest struct {
	Body string `json:"body"`
}

type userPostRequest struct {
	Email string `json:"email"`
}

type succesfullChirpPostResponse struct {
	CleanedBody string `json:"cleaned_body"`
}