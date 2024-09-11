package database

import (
	"os"
	"fmt"
	"encoding/json"
	"sort"
	"golang.org/x/crypto/bcrypt"
	"github.com/niccolot/Chirpy/internal/errors"
)


func (db *DB) CreateChirp(body string) (Chirp, *errors.CodedError) {
	err := validateChirp(&body)
	if err != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to validate chirp: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return Chirp{}, &e
	}

	numChirps, errSize := db.GetNumChirps()
	if errSize != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to get database size: %w, function: %s", errSize, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return Chirp{}, &e
	}

	chirp := Chirp{
		Body: body,
		Id: numChirps+1,
	}

	return chirp, nil
}

func (db *DB) CreateUser(email string, password string) (User, *errors.CodedError) {
	numUsers, errSize := db.GetNumUsers()
	if errSize != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to get database size: %w, function: %s", errSize, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return User{}, &e
	}

	password_bytes := []byte(password)
	hash, errHashing := bcrypt.GenerateFromPassword(password_bytes, bcrypt.DefaultCost)
	if errHashing != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("error hashing password: %w, function: %s", errHashing, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return User{}, &e
	}

	user := User{
		Email: email,
		Password: string(hash),
		Id: numUsers+1,
	}

	return user, nil
}

func (db *DB) UpdateUser(userId int, email string, password string) (Updateduser, *errors.CodedError) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed to load database file: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return Updateduser{}, &e
	}

	//dbStruct.mux.RLock()
	//defer dbStruct.mux.Unlock()
	user := dbStruct.Users[userId]
	user.Email = email
	password_bytes := []byte(password)
	hash, errHashing := bcrypt.GenerateFromPassword(password_bytes, bcrypt.DefaultCost)
	if errHashing != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("error hashing password: %w, function: %s", errHashing, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		}
		return Updateduser{}, &e
	}
	
	user.Password = string(hash)
	dbStruct.Users[userId] = user
	db.WriteDB(&dbStruct)

	updateduser := Updateduser{
		Email: dbStruct.Users[userId].Email,
		Id: userId,
	}

	return updateduser, nil
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

func (db *DB) SearchUserEmail(email string) (bool, int, *errors.CodedError) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		e := errors.CodedError{
			Message: fmt.Errorf("failed loading database: %w, function: %s", err, errors.GetFunctionName()).Error(),
			StatusCode: 500,
		} 
		return false, 0, &e
	}

	found := false
	var userIdx int
	for i, user := range(dbStruct.Users) {
		if user.Email == email {
			found = true
			userIdx = i
			break
		}
	}

	return found, userIdx, nil
}