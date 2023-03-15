package database

import (
	"fmt"
	"sync"
)

// InMemoryDatabase Json file as a database to store key-value pairs
type InMemoryDatabase struct {
	m       sync.RWMutex
	dataMap map[string]string
}

// Init Initialise the database
func (db *InMemoryDatabase) Init() error {
	db.dataMap = make(map[string]string)
	return nil
}

// Get reads json file and returns a map containing data
func (db *InMemoryDatabase) Get(key string) (string, error) {
	db.m.RLock()
	defer db.m.RUnlock()
	if _, ok := db.dataMap[key]; !ok {
		return "", fmt.Errorf("database entry not found")
	}
	return db.dataMap[key], nil
}

// Update updates the db with key/value pair
func (db *InMemoryDatabase) Update(key, value string) error {
	db.m.Lock()
	defer db.m.Unlock()

	db.dataMap[key] = value
	return nil
}

// Delete deletes the database entry
func (db *InMemoryDatabase) Delete(key string) error {
	db.m.Lock()
	defer db.m.Unlock()

	// delete the entry
	delete(db.dataMap, key)
	return nil
}

// Contains Checks for the presence for an entry in the database
func (db *InMemoryDatabase) Contains(key string) bool {
	db.m.RLock()
	defer db.m.RUnlock()
	if _, ok := db.dataMap[key]; !ok {
		return false
	}
	return true
}
