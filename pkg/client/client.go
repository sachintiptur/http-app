package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Data struct {
	Key, Val string
}

// Function to create http request based on the HTTP method
func createHTTPRequest(method string, userData Data) (*http.Request, error) {
	var url string
	url = "http://localhost:8080"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		_ = fmt.Errorf("create request failed: %v", err)
	}

	query := req.URL.Query()

	switch method {
	case http.MethodGet:
		query.Add("key", userData.Key)

	case http.MethodPut:
		query.Add("key", userData.Key)
		query.Add("value", userData.Val)

	case http.MethodDelete:
		query.Add("key", userData.Key)

	}
	req.URL.RawQuery = query.Encode()

	return req, err

}

// Function to send HTTP request and parse the response
func PrepareAndSendHTTPRequest(method string, tmp Data) (string, error) {
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
