package jsondatabase

import (
	"sync"
)


type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
	AuthorId int `json:"author_id"`
}

type DB struct {
	path string
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
	Mux *sync.RWMutex
}

type User struct {
	Id int `json:"id"`
	Email string `json:"email"` 
	Password string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

type Updateduser struct {
	Id int `json:"id"`
	Email string `json:"email"` 
	IsChirpyRed bool `json:"is_chirpy_red"`
}

