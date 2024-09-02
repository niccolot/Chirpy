package database

import (
	"sync"
	"os"
	"fmt"
	"encoding/json"
	"sort"
	"github.com/niccolot/Chirpy/internal/errors"
)


func NewDB(path string) (*DB, *errors.CodedError) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			e := errors.CodedError{
				Message:   fmt.Errorf("failed to create database file %w, function: %s", err, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			return nil, &e
		}
		defer file.Close()
		fmt.Println("Database file created:", path)
	} 

	db := &DB{
		path: path,
		mux: &sync.RWMutex{},
	}

	return db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, *errors.CodedError) {
	err := validateChirp(&body)
	if err != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("failed to validate chirp: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return Chirp{}, &e
	}

	chirps, errSize := db.GetNumChirps()
	if errSize != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("failed to get database size: %w, function: %s", errSize, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return Chirp{}, &e
	}

	c := Chirp{
		Body: body,
		Id: chirps+1,
	}

	return c, nil
}

func (db *DB) CreateUser(email string) (User, *errors.CodedError) {
	chirps, errSize := db.GetNumUsers()
	if errSize != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("failed to get database size: %w, function: %s", errSize, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return User{}, &e
	}

	u := User{
		Email: email,
		Id: chirps+1,
	}

	return u, nil
}

func (db *DB) GetChirps() ([]Chirp, *errors.CodedError) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to load database file: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return nil, &e
	}

	len, err := db.GetNumChirps()
	if err != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to get database size: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return nil, &e
	}

	chirpsSlice := make([]Chirp, len)
	for i, chirp := range dbStruct.Chirps {
		chirpsSlice[i-1] = chirp
	}

	sort.Slice(chirpsSlice, func(i, j int) bool {
		return chirpsSlice[i].Id < chirpsSlice[j].Id 
	})

	return chirpsSlice, nil
}

func (db *DB) GetChirpID(id int) (Chirp, *errors.CodedError) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to load database file: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return Chirp{}, &e
	}

	chirp, ok := dbStruct.Chirps[id]
	if !ok {
		e := errors.CodedError{
			Message: "Chirp not found",
			StatusCode: 404,
		}
		return Chirp{}, &e
	}

	return chirp, nil
}

func (db *DB) GetNumChirps() (int, *errors.CodedError) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("error reading database file: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return 0, &e
	}

	return len(dbStruct.Chirps), nil
}

func (db *DB) GetNumUsers() (int, *errors.CodedError) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("error reading database file: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return 0, &e
	}

	return len(dbStruct.Users), nil
}

func (db *DB) LoadDB() (DBStructure, *errors.CodedError) {
	fileContent, err := os.ReadFile(db.path)
	if err != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("error reading database file: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return GetDBStruct(), &e
	}

	if len(fileContent) == 0 {
		return GetDBStruct(), nil
	}

	dbStruct := GetDBStruct()

	err = json.Unmarshal(fileContent, &dbStruct)
	if err != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("error unmarshaling json: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return GetDBStruct(), &e
	}

	return dbStruct, nil
}

func (db *DB) ensureDB() *errors.CodedError {
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			e := errors.CodedError{
				Message:   fmt.Errorf("failed to create databse file: %w, function: %s", err, errors.GetFunctionName()).Error(),
				StatusCode: 500,
			}
			return &e
		}
		defer file.Close()
		fmt.Println("Database file created:", db.path)
	} 

	return nil
}

func (db *DB) WriteDB(dbStructure *DBStructure) *errors.CodedError {
	err := db.ensureDB()
	if err != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("%w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return &e
	}

	jsonData, errMarshal := json.MarshalIndent(dbStructure, "", "  ")
	if errMarshal != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to marshal map: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return &e
	}

	cwd , errGetwd := os.Getwd()
	if errGetwd != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to get working directory: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return &e
	} 

	path := cwd + "/" + db.path
	dbFile, errOpen := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if errOpen != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to open database file: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return &e
	}

	defer dbFile.Close()
	_, errWrite := dbFile.Write(jsonData)
	if errWrite != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to write to file: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return &e
	}

	return nil
}

func GetDBStruct() DBStructure {
	return DBStructure{
		make(map[int]Chirp),
		make(map[int]User),
	}
}