package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	database "github.com/sachintiptur/http-app/util"
)

type dbStruct struct {
	dbIntf database.Database
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
func processGET(dbI database.Database, resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if ok := dbI.Contains(key); !ok {
		http.Error(resp, "Database entry not found", http.StatusNotFound)
	} else {
		//data, ok := dbI.(*database.MemData)
		data, _ := dbI.Read()
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintf("Data found for key %s: %s", key, data[key])))
	}
	return
}

// Process the HTTP PUT request
func processPUT(dbI database.Database, resp http.ResponseWriter, req *http.Request) {

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
func processDELETE(dbI database.Database, resp http.ResponseWriter, req *http.Request) {
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
func (dbS dbStruct) processHTTPRequests(resp http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	if status, err := ValidateData(key, val); err != nil {
		http.Error(resp, err.Error(), status)
		return
	}
	// Handle http methods
	switch req.Method {
	case http.MethodGet:
		processGET(dbS.dbIntf, resp, req)
	case http.MethodPut:
		processPUT(dbS.dbIntf, resp, req)
	case http.MethodDelete:
		processDELETE(dbS.dbIntf, resp, req)
	default:
		http.Error(resp, "Invalid method or not handled", http.StatusMethodNotAllowed)
	}

}

// main function takes server address as input
// and starts listening for http requests
func main() {
	addr := flag.String("addr", ":8080", "Server address string")
	dbType := flag.String("db", "map", "Database to use, supported values are [map, json]")
	flag.Parse()

	// map of supported database types
	var db = map[string]database.Database{"map": &database.JsonData{}, "value": &database.MemData{}}
	var dbS dbStruct
	dbS.dbIntf = db[*dbType]
	dbS.dbIntf.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/", dbS.processHTTPRequests)

	mwMux := NewLogInfo(mux)
	log.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(*addr, mwMux))

}
