package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// Json file as a database to store key-value pairs
type JsonData struct {
	m       sync.RWMutex
	dataMap map[string]string
}

// Database interface
type Database interface {
	Init() error
	Read() (map[string]string, error)
	Write(map[string]string) error
	Contains(key string) bool
}

// DB file name
var JSON_DB = "kvstore.json"

const (
	MAX_DB_ENTRY = 1000
	MAX_KEY_SZ   = 16
	MAX_VAL_SZ   = 32
)

// Initialise the database
// Initialise the data map and creates json file
func (db *JsonData) Init() error {
	db.dataMap = make(map[string]string)
	file, err := os.Create(JSON_DB)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("creating json file failed: %s", err)
	}
	return nil
}

// Database read function
// reads json file and returns a map containing data
func (db *JsonData) Read() (map[string]string, error) {

	db.m.RLock()
	defer db.m.RUnlock()

	file, err := os.OpenFile(JSON_DB, os.O_RDWR, 0655)
	if err != nil {
		return nil, fmt.Errorf("Opening json file failed: %s", err)
	}
	defer file.Close()

	// read json file
	data, _ := ioutil.ReadAll(file)
	if len(data) == 0 {
		return db.dataMap, nil
	}
	// unmarshal json data to data map
	err = json.Unmarshal([]byte(data), &db.dataMap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %s\n", err)
	}

	return db.dataMap, nil
}

// Database write function
// writes the data map into json file
func (db *JsonData) Write(data map[string]string) error {
	db.m.Lock()
	defer db.m.Unlock()

	db.dataMap = data
	file, err := os.OpenFile(JSON_DB, os.O_RDWR, 0655)
	if err != nil {
		return fmt.Errorf("Opening json file failed: %s", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(db.dataMap)
	if err != nil {
		return fmt.Errorf("could not marshal json: %s\n", err)

	}
	err = ioutil.WriteFile(file.Name(), jsonData, 0655)

	if err != nil {
		return fmt.Errorf("Writing record to json file failed: %s", err)
	}

	return nil

}

// Contains function
// Checks for the presence for entry in the database
func (db *JsonData) Contains(key string) bool {
	db.Read()
	if _, ok := db.dataMap[key]; !ok {
		return false
	}
	return true
}

// Map as a database to store key-value pairs
type MemData JsonData

func (db *MemData) Init() error {
	db.dataMap = make(map[string]string)
	return nil
}

func (db *MemData) Read() (map[string]string, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	return db.dataMap, nil
}

func (db *MemData) Write(data map[string]string) error {
	db.m.Lock()
	defer db.m.Unlock()

	db.dataMap = data
	return nil
}

func (db *MemData) Contains(key string) bool {
	db.m.Lock()
	defer db.m.Unlock()
	if _, ok := db.dataMap[key]; !ok {
		return false
	}
	return true
}
