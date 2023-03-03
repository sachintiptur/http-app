package main

import (
	"log"
	"net/http"
)

var database = make(map[string]string)

type httpRequest http.Request

// type processHTTP interface {
// 	processGET(resp http.ResponseWriter)
// 	processPUT(resp http.ResponseWriter)
// 	processDELETE(resp http.ResponseWriter)
// }

// func (req *httpRequest) processGET(resp http.ResponseWriter) {
// 	key := req.URL.Query().Get("key")

// 	if _, ok := database[key]; !ok {
// 		resp.WriteHeader(http.StatusNotFound)
// 	} else {
// 		resp.WriteHeader(http.StatusOK)
// 		resp.Write([]byte(database[key]))
// 	}

// }

func processGET(resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if _, ok := database[key]; !ok {
		resp.WriteHeader(http.StatusNotFound)
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(database[key]))
	}

}

func processPUT(resp http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get("key")
	val := req.URL.Query().Get("value")

	if _, ok := database[key]; !ok {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("Database updated with new key/value pair"))
	} else {
		resp.WriteHeader(http.StatusFound)
		resp.Write([]byte("Updated the existing key/vlaue pair"))
	}

	database[req.URL.Query().Get("key")] = val

}

func processDELETE(resp http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	if _, ok := database[key]; ok {
		delete(database, key)
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("Database entry deleted"))

	} else {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("Database entry not found"))
	}

}

func processHTTPRequests(resp http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("Requested URL not found"))
		http.Error(resp, "not found.", http.StatusNotFound)
		return
	}

	var key, val string
	key = req.URL.Query().Get("key")
	val = req.URL.Query().Get("value")
	log.Println(key, val)

	switch req.Method {

	case http.MethodGet:
		//var r *httpRequest
		// r = (*httpRequest)(req)
		processGET(resp, req)
		// if _, ok := database[key]; !ok {
		// 	resp.WriteHeader(http.StatusNotFound)
		// } else {
		// 	resp.WriteHeader(http.StatusOK)
		// 	resp.Write([]byte(database[key]))
		// }

	case http.MethodPut:

		if _, ok := database[key]; !ok {
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte("Database updated with new key/value pair"))
		} else {
			resp.WriteHeader(http.StatusFound)
			resp.Write([]byte("Updated the existing key/vlaue pair"))
		}

		database[req.URL.Query().Get("key")] = req.URL.Query().Get("value")

	case http.MethodDelete:
		if _, ok := database[key]; ok {
			delete(database, key)
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte("Database entry deleted"))

		} else {
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte("Database entry not found"))
		}

	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", processHTTPRequests)

	log.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
