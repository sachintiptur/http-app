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

func main() {

	var url string
	method := flag.String("m", "", "HTTP method")
	key := flag.String("k", "", "Key for the data")
	value := flag.String("v", "", "Value of the data")

	flag.Parse()

	client := http.Client{Timeout: time.Duration(1) * time.Second}
	url = "http://localhost:8080"

	req, err := http.NewRequest(*method, url, nil)
	if err != nil {
		_ = fmt.Errorf("create request failed: %v", err)
	}

	query := req.URL.Query()

	switch *method {
	case http.MethodGet:
		query.Add("key", *key)

	case http.MethodPut:
		query.Add("key", *key)
		query.Add("value", *value)

	case http.MethodDelete:
		query.Add("key", *key)

	}

	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http request failed. error %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(resp.Status)
	log.Println(string(body[:]))

}
