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

// ValidateData Validate the data sent by client
// Check for key and value length
func ValidateData(key, val string) error {

	if len(key) > database.MaxKeySz {
		return fmt.Errorf("key length is greater than %d", database.MaxKeySz)
	} else if len(val) > database.MaxValSz {
		return fmt.Errorf("value length is greater than %d", database.MaxValSz)
	} else {
		return nil
	}
}

// LogInfo Logging middle-ware handler struct
type LogInfo struct {
	handler http.Handler
}

func NewLogInfo(reqHandler http.Handler) *LogInfo {
	return &LogInfo{reqHandler}
}

// ServeHTTP Interface implementation for LogInfo
// Wraps the actual http handler functions with the
// logging middleware functions
func (l *LogInfo) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	log.Printf("METHOD: %s PATH: %s KEY: %s", req.Method, req.URL, req.URL.Query().Get("key"))
	l.handler.ServeHTTP(resp, req)
	log.Printf("Time elapsed: %v", time.Since(start))
}

// ProcessGET Process the HTTP GET request
func ProcessGET(dbI database.Database, resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	// read the entry from database
	value, err := dbI.Get(key)
	if err != nil {
		http.Error(resp, "Database entry not found", http.StatusNotFound)
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintf("Data found for key %s: %s", key, value)))
	}
	return
}

// ProcessPUT Process the HTTP PUT request
func ProcessPUT(dbI database.Database, resp http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	// Check if db contains the key, the 'ok' flag is later used
	// to send the appropriate http status
	ok := dbI.Contains(key)

	// write to database
	err := dbI.Update(key, val)
	if err != nil {
		http.Error(resp, "database write failed", http.StatusNotModified)
	} else {
		if !ok {
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte("Database updated with new key/value pair"))
		} else {
			// patch request case
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte("Updated the existing key/value pair"))
		}
	}

}

// ProcessDELETE Process the HTTP DELETE request
func ProcessDELETE(dbI database.Database, resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if ok := dbI.Contains(key); !ok {
		http.Error(resp, "Database entry not found", http.StatusNotFound)
		return
	}

	if err := dbI.Delete(key); err != nil {
		http.Error(resp, err.Error(), http.StatusNotModified)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Database entry deleted"))
	return

}

// ProcessHTTPRequests Handler function to handle HTTP requests
func (dbS DbStruct) ProcessHTTPRequests(resp http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	if err := ValidateData(key, val); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
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
