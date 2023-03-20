package server

import (
	"fmt"
	"log"
	"net/http"

	database "github.com/sachintiptur/http-app/pkg/database"
)

// DatabaseHandler holds the database interface type
type DatabaseHandler struct {
	Db database.Database
}

// ValidateData Validate the data sent by client
// Check for key and value length
func ValidateData(key, val string) error {
	if len(key) > database.MaxKeySz {
		return fmt.Errorf("key length is greater than %d", database.MaxKeySz)
	} else if len(val) > database.MaxValSz {
		return fmt.Errorf("value length is greater than %d", database.MaxValSz)
	}
	return nil
}

// ProcessGET Process the HTTP GET request
func (dbHandler *DatabaseHandler) ProcessGET(resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	// read the entry from database
	value, err := dbHandler.Db.Get(key)
	if err != nil {
		http.Error(resp, fmt.Sprintf("Data not found for key: %s", key), http.StatusNotFound)
	} else {
		resp.WriteHeader(http.StatusOK)
		_, err = resp.Write([]byte(fmt.Sprintf("Data found for key: %s,%s", key, value)))
		if err != nil {
			return
		}
	}
}

// ProcessPUT Process the HTTP PUT request
func (dbHandler *DatabaseHandler) ProcessPUT(resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	// Check if db contains the key, the 'ok' flag is later used
	// to send the appropriate http status
	ok := dbHandler.Db.Contains(key)

	// write to database
	err := dbHandler.Db.Update(key, val)
	if err != nil {
		http.Error(resp, "database write failed", http.StatusNotModified)
	} else {
		if !ok {
			resp.WriteHeader(http.StatusCreated)
			_, err = resp.Write([]byte(fmt.Sprintf("Database updated with new key/value pair: %s,%s", key, val)))
			if err != nil {
				log.Println("writing response failed")
				return
			}
		} else {
			// patch request case
			resp.WriteHeader(http.StatusOK)
			_, err = resp.Write([]byte(fmt.Sprintf("Updated the existing key/value pair: %s,%s", key, val)))
			if err != nil {
				log.Println("writing response failed")
				return
			}
		}
	}
}

// ProcessDELETE Process the HTTP DELETE request
func (dbHandler *DatabaseHandler) ProcessDELETE(resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if ok := dbHandler.Db.Contains(key); !ok {
		http.Error(resp, fmt.Sprintf("Database entry not found for key: %s", key), http.StatusNotFound)
		return
	}

	if err := dbHandler.Db.Delete(key); err != nil {
		http.Error(resp, err.Error(), http.StatusNotModified)
		return
	}

	resp.WriteHeader(http.StatusOK)
	_, err := resp.Write([]byte(fmt.Sprintf("Database entry deleted for key: %s", key)))
	if err != nil {
		log.Println("writing response failed")
	}
}

// ProcessHTTPRequests Handler function to handle HTTP requests
func (dbHandler *DatabaseHandler) ProcessHTTPRequests(resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	if err := ValidateData(key, val); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	// Handle http methods
	switch req.Method {
	case http.MethodGet:
		dbHandler.ProcessGET(resp, req)
	case http.MethodPut:
		dbHandler.ProcessPUT(resp, req)
	case http.MethodDelete:
		dbHandler.ProcessDELETE(resp, req)
	default:
		http.Error(resp, "Invalid method or not handled", http.StatusMethodNotAllowed)
	}
}
