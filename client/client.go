package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type data struct {
	key, val string
}

// Function to create http request based on the HTTP method
func createHTTPRequest(method string, userData data) (*http.Request, error) {
	var url string
	url = "http://localhost:8080"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		_ = fmt.Errorf("create request failed: %v", err)
	}

	query := req.URL.Query()

	switch method {
	case http.MethodGet:
		query.Add("key", userData.key)

	case http.MethodPut:
		query.Add("key", userData.key)
		query.Add("value", userData.val)

	case http.MethodDelete:
		query.Add("key", userData.key)

	}
	req.URL.RawQuery = query.Encode()

	return req, err

}

// Function to send HTTP request and parse the response
func PrepareAndSendHTTPRequest(method string, tmp data) (string, error) {
	client := http.Client{Timeout: time.Duration(1) * time.Second}

	req, err := createHTTPRequest(method, tmp)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http request failed. error %s", err)
		return "HTTP request failed", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil

}

// Main function takes http method, key and value as the command
// line arguments. Build and send appropriate HTTP request to
// the server.
func main() {

	method := flag.String("m", "", "HTTP method")
	key := flag.String("k", "", "Key for the data")
	value := flag.String("v", "", "Value of the data")

	flag.Parse()

	var tmp = data{key: *key, val: *value}

	resp, err := PrepareAndSendHTTPRequest(*method, tmp)
	if err != nil {
		fmt.Printf("PrepareAndSendHTTPRequest failed with error %s", err)
		return
	}

	log.Println(resp)

}
