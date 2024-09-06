package database

import (
	"sync"
)


type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

type User struct {
	Id int `json:"id"`
	Email string `json:"email"` 
	Password string `json:"password"`
}