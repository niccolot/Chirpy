package database

import (
	"strings"
	"github.com/niccolot/Chirpy/internal/errors"
)


func validateChirp(body *string) error {
	maxChirpLength := 140
	if len(*body) > maxChirpLength {
		e := errors.CodedError{
			Message:   "Error: chirp is too long",
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
