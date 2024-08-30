package database

import (
	"sync"
	"os"
	"fmt"
	"encoding/json"
	"sort"
	"github.com/niccolot/Chirpy/internal/errors"
)


func NewDB(path string) (*DB, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			e := errors.CodedError{
				Message:   fmt.Errorf("%w", err).Error(),
				StatusCode: 500,
			}
			fmt.Println("error creating database")
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

func (db *DB) CreateChirp(body string) (Chirp, error) {
	err := validateChirp(&body)
	if err != nil {
		return Chirp{}, err
	}

	chirps, err := db.GetDBLength()
	if err != nil {
		return Chirp{}, err
	}

	c := Chirp{
		Body: body,
		Id: string(chirps+1),
	}

	return c, nil
}

func (db *DB) GetDBLength() (int, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return 0, err
	}

	return len(dbStruct.Chirps), nil
}

func (db *DB) LoadDB() (DBStructure, error) {
	fileContent, err := os.ReadFile(db.path)
	if err != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("%w", err).Error(),
			StatusCode: 500,
		}
		fmt.Println("Error reading  database file:", err)
		return GetDBStruct(), &e
	}

	if len(fileContent) == 0 {
		return GetDBStruct(), nil
	}

	dbStruct := GetDBStruct()

	err = json.Unmarshal(fileContent, &dbStruct)
	if err != nil {
		e := errors.CodedError{
			Message:   fmt.Errorf("%w", err).Error(),
			StatusCode: 500,
		}
		fmt.Println("Error parsing JSON:", err)
		return GetDBStruct(), &e
	}

	return dbStruct, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.LoadDB()
	if err != nil {
		return nil, err
	}

	len, err := db.GetDBLength()
	if err != nil {
		return nil, err
	}

	chirpsSlice := make([]Chirp, len)
	for i, chirp := range dbStruct.Chirps {
		chirpsSlice[i] = chirp
	}

	sort.Slice(chirpsSlice, func(i, j int) bool {
		return chirpsSlice[i].Id < chirpsSlice[j].Id 
	})

	return chirpsSlice, nil
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			e := errors.CodedError{
				Message:   fmt.Errorf("%w", err).Error(),
				StatusCode: 500,
			}
			fmt.Println("error creating database")
			return &e
		}
		defer file.Close()
		fmt.Println("Database file created:", db.path)
	} else {
		fmt.Println("Database file already exists:", db.path)
	}

	return nil
}

func (db *DB) WriteDB(dbStructure *DBStructure) error {
	err := db.ensureDB()
	if err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(dbStructure, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal map: %v", err)
	}

	path := "/home/nico/repos/Chirpy/" + db.path
	dbFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	
	if err != nil {
		return err
	}

	defer dbFile.Close()
	_, err = dbFile.Write(jsonData)
	if err != nil {
		fmt.Println("pippo")

		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func GetDBStruct() DBStructure {
	return DBStructure{
		make(map[int]Chirp),
	}
}