package database

import (
	"encoding/json"
	"fmt"
	"os"
)

// JsonDatabase Map as a database to store key-value pairs
type JsonDatabase InMemoryDatabase

// Init Initialise the database
// Initialises the data map and creates json file
func (db *JsonDatabase) Init() error {
	db.dataMap = make(map[string]string)
	file, err := os.Create(JsonDB)
	if err != nil {
		return fmt.Errorf("creating json file failed: %s", err)
	}
	defer file.Close()

	// initialise with empty data
	jsonData, err := json.Marshal(db.dataMap)
	if err != nil {
		fmt.Println("json marshal failed")
		return err

	}
	err = os.WriteFile(JsonDB, jsonData, 0o655)

	if err != nil {
		fmt.Println("writing to file failed")
		return err
	}
	return err
}

// Get reads json file and
// returns a map containing data
func (db *JsonDatabase) Get(key string) (val string, err error) {
	db.m.RLock()
	defer db.m.RUnlock()

	data, err := os.ReadFile(JsonDB)
	if err != nil {
		fmt.Println("reading json file failed")
		return "", err
	}

	err = json.Unmarshal(data, &db.dataMap)
	if err != nil {
		fmt.Println(" unmarshalling json failed")
		return "", err
	}
	if _, ok := db.dataMap[key]; !ok {
		return "", fmt.Errorf("key not present")
	}

	return db.dataMap[key], nil
}

// Update writes the data map into json file
func (db *JsonDatabase) Update(key, value string) error {
	db.m.Lock()
	defer db.m.Unlock()

	data, err := os.ReadFile(JsonDB)
	if err != nil {
		fmt.Println("reading json file failed")
		return err
	}
	err = json.Unmarshal(data, &db.dataMap)
	if err != nil {
		fmt.Println(" unmarshalling json data failed")
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
	err = os.WriteFile(JsonDB, jsonData, 0o655)

	if err != nil {
		fmt.Println(" writing to json file failed")
		return err
	}
	return nil
}

// Delete deletes the database entry
func (db *JsonDatabase) Delete(key string) error {
	db.m.Lock()
	defer db.m.Unlock()
	data, err := os.ReadFile(JsonDB)
	if err != nil {
		fmt.Println("reading json file failed")
		return err
	}

	err = json.Unmarshal(data, &db.dataMap)
	if err != nil {
		return fmt.Errorf("could not unmarshal json: %s\n", err)
	}

	delete(db.dataMap, key)
	jsonData, err := json.Marshal(db.dataMap)
	if err != nil {
		return fmt.Errorf("could not marshal json: %s\n", err)
	}
	err = os.WriteFile(JsonDB, jsonData, 0o655)

	if err != nil {
		return err
	}
	return nil
}

// Contains Checks for the presence for an entry in the database
func (db *JsonDatabase) Contains(key string) bool {
	if _, err := db.Get(key); err != nil {
		return false
	}
	return true
}
