package jsondatabase

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"github.com/niccolot/Chirpy/internal/errors"
)


func NewDB(path string) (*DB, *errors.CodedError) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			e := errors.CodedError{
				Message: fmt.Errorf("failed to create database file %w, function: %s", err, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			return nil, &e
		}
		defer file.Close()
		fmt.Printf("Database file created: %s\n", path)
	} 

	db := &DB{
		path: path,
	}

	return db, nil
}

func validateChirp(body *string) error {
	maxChirpLength := 140
	if len(*body) > maxChirpLength {
		e := errors.CodedError{
			Message:   "Error: chirp is too long\n",
			StatusCode: 400,
		}
		return &e
	}

	cleanProfanity(body)
	return nil
}

func cleanProfanity(body *string) {
	
	badWords := map[string]bool{		
		"kerfuffle": true,
		"sharbert": true,
		"fornax": true,
	}

	censor := "****"
	words := strings.Fields(*body)

	for i, word := range  words {
		if badWords[strings.ToLower(word)] {
			words[i] = censor
		}
	}

	*body = strings.Join(words, " ")
}

func GetDBStruct() DBStructure {
	return DBStructure{
		make(map[int]Chirp),
		make(map[int]User),
		&sync.RWMutex{},
	}
}