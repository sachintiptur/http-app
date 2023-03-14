package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// JsonData Json file as a database to store key-value pairs
type JsonData struct {
	m       sync.RWMutex
	dataMap map[string]string
}

// Database interface
type Database interface {
	Init() error
	Get(key string) (string, error)
	Update(key, value string) error
	Delete(key string) error
	Contains(key string) bool
}

// JsonDB database file name
var JsonDB = "kvstore.json"

const (
	MaxDbEntry = 1000
	MaxKeySz   = 16
	MaxValSz   = 32
)

// Init Initialise the database
// Initialises the data map and creates json file
func (db *JsonData) Init() error {
	db.dataMap = make(map[string]string)
	file, err := os.Create(JsonDB)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("creating json file failed: %s", err)
	}

	// initialise with empty data
	jsonData, err := json.Marshal(db.dataMap)
	if err != nil {
		fmt.Print("json marshal failed")
		return err

	}
	err = os.WriteFile(JsonDB, jsonData, 0655)

	if err != nil {
		fmt.Print("writing to file failed")
		return err
	}
	return err
}

// Get function reads json file and
// returns a map containing data
func (db *JsonData) Get(key string) (val string, err error) {

	db.m.RLock()
	defer db.m.RUnlock()

	data, err := os.ReadFile(JsonDB)
	if err != nil {
		fmt.Print("reading json file failed")
		return "", err
	}

	// unmarshal json data to data map
	err = json.Unmarshal([]byte(data), &db.dataMap)
	if err != nil {
		fmt.Print(" unmarshalling json failed")
		return "", err
	}
	if _, ok := db.dataMap[key]; !ok {
		return "", fmt.Errorf("key not present")
	}

	return db.dataMap[key], nil
}

// Update function
// writes the data map into json file
func (db *JsonData) Update(key, value string) error {
	db.m.Lock()
	defer db.m.Unlock()

	data, err := os.ReadFile(JsonDB)
	if err != nil {
		fmt.Print("reading json file failed")
		return err
	}
	err = json.Unmarshal([]byte(data), &db.dataMap)
	if err != nil {
		fmt.Print(" unmarshalling json data failed")
		return err
	}

	if len(db.dataMap) > MaxDbEntry {
		return fmt.Errorf("database maximum size %d reached", MaxDbEntry)
	}
	db.dataMap[key] = value

	jsonData, err := json.Marshal(db.dataMap)
	if err != nil {
		return fmt.Errorf("could not marshal json: %s\n", err)

	}
	err = os.WriteFile(JsonDB, jsonData, 0655)

	if err != nil {
		fmt.Print(" writing to json file failed")
		return err
	}
	return nil
}

func (db *JsonData) Delete(key string) error {

	db.m.Lock()
	defer db.m.Unlock()
	data, err := os.ReadFile(JsonDB)
	if err != nil {
		fmt.Print("reading json file failed")
		return err
	}

	// unmarshal json data to data map
	err = json.Unmarshal([]byte(data), &db.dataMap)
	if err != nil {
		return fmt.Errorf("could not unmarshal json: %s\n", err)
	}
	// delete the entry
	delete(db.dataMap, key)
	jsonData, err := json.Marshal(db.dataMap)
	if err != nil {
		return fmt.Errorf("could not marshal json: %s\n", err)

	}
	err = os.WriteFile(JsonDB, jsonData, 0655)

	if err != nil {
		return err
	}
	return nil

}

// Contains function
// Checks for the presence for entry in the database
func (db *JsonData) Contains(key string) bool {
	if _, err := db.Get(key); err != nil {
		return false
	}
	return true
}

// MemData Map as a database to store key-value pairs
type MemData JsonData

// Init Initialise the database
func (db *MemData) Init() error {
	db.dataMap = make(map[string]string)
	return nil
}

// Get Read function reads json file and returns a map containing data
func (db *MemData) Get(key string) (string, error) {
	db.m.RLock()
	defer db.m.RUnlock()
	if _, ok := db.dataMap[key]; !ok {
		return "", fmt.Errorf("database entry not found")
	}
	return db.dataMap[key], nil
}

// Update function updates the db with key/value pair
func (db *MemData) Update(key, value string) error {
	db.m.Lock()
	defer db.m.Unlock()

	db.dataMap[key] = value
	return nil
}

// Delete function deletes the database entry
func (db *MemData) Delete(key string) error {
	db.m.Lock()
	defer db.m.Unlock()

	// delete the entry
	delete(db.dataMap, key)
	return nil
}

// Contains function
// Checks for the presence for entry in the database
func (db *MemData) Contains(key string) bool {
	db.m.RLock()
	defer db.m.RUnlock()
	if _, ok := db.dataMap[key]; !ok {
		return false
	}
	return true
}
