package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Data struct holding user input data
type Data struct {
	Key, Val string
}

// createHTTPRequest create http request based on the HTTP method
func createHTTPRequest(method string, userData Data) (*http.Request, error) {
	url := "http://localhost:8080"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return req, fmt.Errorf("create request failed: %v", err)
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

// SendHTTPRequest sends HTTP requests and parse the response
func SendHTTPRequest(method string, data Data) (string, error) {
	client := http.Client{Timeout: time.Duration(1) * time.Second}

	req, err := createHTTPRequest(method, data)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
