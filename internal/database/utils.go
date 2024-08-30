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

/*
func validateChirpGOLD(w http.ResponseWriter, r *http.Request) {
	maxChirpLength := 140

	decoder := json.NewDecoder(r.Body)
	res := response{}
	err := decoder.Decode(&res)

	if err != nil {
		respondWithErrorGOLD(w, err)
		return 
	}

	if len(res.Body) > maxChirpLength {
		respondTooLongChirpGOLD(w)
		return    
	}

	cleanProfanity(&res.Body)
	respondWithSuccessGOLD(w, &res)
}

func respondTooLongChirpGOLD(w http.ResponseWriter) {
	w.WriteHeader(400)
	errResp := errResponse{
		Error: "Chirp is too long",
	}

	dat, err := json.Marshal(errResp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func respondWithSuccessGOLD(w http.ResponseWriter, r *response) {
	succResp := succesfullResponse{
		CleanedBody: r.Body,
	}
	
	dat, err := json.Marshal(succResp)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}
*/

