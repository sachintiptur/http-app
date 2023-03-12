package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	database "github.com/sachintiptur/http-app/pkg/util"
)

type DbStruct struct {
	DbIntf database.Database
}

var data map[string]string

// Validate the data sent by client
// Check for key and value length
func ValidateData(key, val string) (int, error) {

	if len(key) > database.MAX_KEY_SZ {
		return http.StatusBadRequest, fmt.Errorf("Key length is greater than %d", database.MAX_KEY_SZ)
	} else if len(val) > database.MAX_VAL_SZ {
		return http.StatusBadRequest, fmt.Errorf("Value length is greater than %d", database.MAX_VAL_SZ)
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
	log.Printf("METHOD: %s PATH: %s KEY: %s", req.Method, req.URL, req.URL.Query().Get("key"))
	l.handler.ServeHTTP(resp, req)
	log.Printf("Time elapsed: %v", time.Since(start))
}

// Process the HTTP GET request
func ProcessGET(dbI database.Database, resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if ok := dbI.Contains(key); !ok {
		http.Error(resp, "Database entry not found", http.StatusNotFound)
	} else {
		data, _ := dbI.Read()
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintf("Data found for key %s: %s", key, data[key])))
	}
	return
}

// Process the HTTP PUT request
func ProcessPUT(dbI database.Database, resp http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	data, err := dbI.Read()
	if err != nil {
		http.Error(resp, "database read failed", http.StatusNotFound)
		return
	}
	if data == nil {
		data = make(map[string]string)
	}
	// check for database limit before writing
	if len(data) == database.MAX_DB_ENTRY {
		http.Error(resp, "Database limit reached", http.StatusInsufficientStorage)
		return
	}

	if _, ok := data[key]; !ok {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("Database updated with new key/value pair"))
	} else {
		http.Error(resp, "Updated the existing key/vlaue pair", http.StatusFound)
	}

	// write to database
	data[key] = val
	err = dbI.Write(data)
	if err != nil {
		http.Error(resp, "database write failed", http.StatusNotModified)
	}

}

// Process the HTTP DELETE request
func ProcessDELETE(dbI database.Database, resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	data, err := dbI.Read()
	if err != nil {
		http.Error(resp, "database read failed", http.StatusNotFound)
		return
	}
	if _, ok := data[key]; !ok {
		http.Error(resp, "Database entry not found", http.StatusNotFound)
		return

	}
	delete(data, key)

	err = dbI.Write(data)
	if err != nil {
		http.Error(resp, "database write failed", http.StatusNotModified)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Database entry deleted"))
	return

}

// Handler function to handle HTTP requests
func (dbS DbStruct) ProcessHTTPRequests(resp http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	if status, err := ValidateData(key, val); err != nil {
		http.Error(resp, err.Error(), status)
		return
	}
	// Handle http methods
	switch req.Method {
	case http.MethodGet:
		ProcessGET(dbS.DbIntf, resp, req)
	case http.MethodPut:
		ProcessPUT(dbS.DbIntf, resp, req)
	case http.MethodDelete:
		ProcessDELETE(dbS.DbIntf, resp, req)
	default:
		http.Error(resp, "Invalid method or not handled", http.StatusMethodNotAllowed)
	}

}
