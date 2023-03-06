package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Map as a database to store key-value pairs
type Database struct {
	m       sync.RWMutex
	dataMap map[string]string
}

const (
	MAX_DB_ENTRY = 1000
	MAX_KEY_SZ   = 16
	MAX_VAL_SZ   = 32
)

var db Database

func (db *Database) InitDatabase() {
	db.dataMap = make(map[string]string)
}

func (db *Database) ValidateData(key, val string) (int, error) {

	if len(db.dataMap) == MAX_DB_ENTRY {
		return http.StatusBadRequest, fmt.Errorf("Database limit of %d reached", MAX_DB_ENTRY)
	} else if len(key) > MAX_KEY_SZ {
		return http.StatusBadRequest, fmt.Errorf("Key length is greater than %d", MAX_KEY_SZ)
	} else if len(val) > MAX_VAL_SZ {
		return http.StatusBadRequest, fmt.Errorf("Value length is greater than %d", MAX_VAL_SZ)
	} else {
		return http.StatusOK, nil
	}
}

// Logging middle-ware handler struct
type LogInfo struct {
	handler http.Handler
}

func NewLogInfo(reqHandler http.Handler) *LogInfo {
	return &LogInfo{reqHandler}
}

// Interface implementation for LogInfo
// Wraps the actual http handler functions with the
// logging middleware functions
func (l *LogInfo) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	log.Printf("METHOD: %s KEY: %s", req.Method, req.URL.Query().Get("key"))
	l.handler.ServeHTTP(resp, req)
	log.Printf("Time elapsed: %v", time.Since(start))
}

// Process the HTTP GET request
func processGET(resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	db.m.RLock()
	defer db.m.RUnlock()

	if _, ok := db.dataMap[key]; !ok {
		http.Error(resp, "Database entry not found", http.StatusNotFound)
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(db.dataMap[key]))
	}

}

// Process the HTTP PUT request
func processPUT(resp http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")
	db.m.Lock()
	defer db.m.Unlock()

	if _, ok := db.dataMap[key]; !ok {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("Database updated with new key/value pair"))
	} else {
		http.Error(resp, "Updated the existing key/vlaue pair", http.StatusNotFound)
	}

	db.dataMap[key] = val

}

// Process the HTTP DELETE request
func processDELETE(resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	db.m.Lock()
	defer db.m.Unlock()

	if _, ok := db.dataMap[key]; ok {
		delete(db.dataMap, key)
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("Database entry deleted"))

	} else {
		http.Error(resp, "Database entry not found", http.StatusNotFound)
	}

}

// Handler function to handle HTTP requests
func processHTTPRequests(resp http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		resp.Write([]byte("Requested URL not found"))
		http.Error(resp, "not found.", http.StatusNotFound)
		return
	}

	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	if status, err := db.ValidateData(key, val); err != nil {
		http.Error(resp, err.Error(), status)
		return
	}

	switch req.Method {
	case http.MethodGet:
		processGET(resp, req)
	case http.MethodPut:
		processPUT(resp, req)
	case http.MethodDelete:
		processDELETE(resp, req)
	default:
		fmt.Printf("Invalid method or not handled")
	}

}

func main() {
	addr := flag.String("addr", ":8080", "Server address string")
	flag.Parse()

	db.InitDatabase()

	mux := http.NewServeMux()
	mux.HandleFunc("/", processHTTPRequests)

	mwMux := NewLogInfo(mux)

	log.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(*addr, mwMux))

}
