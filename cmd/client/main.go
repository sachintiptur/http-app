package main

import (
	"flag"
	"log"

	client "github.com/sachintiptur/http-app/pkg/client"
)

// Client CLI process
// takes http method, key and value as the command
// line arguments. Build and send appropriate HTTP request to
// the server.
func main() {

	method := flag.String("m", "", "HTTP method, supported methods are [GET, PUT, DELETE]")
	key := flag.String("k", "", "Key for the data")
	value := flag.String("v", "", "Value of the data")

	flag.Parse()

	methods := map[string]string{"GET": "", "PUT": "", "DELETE": ""}
	if _, ok := methods[*method]; !ok {
		flag.Usage()
		log.Fatal("method not supported")
	}

	var tmp = client.Data{Key: *key, Val: *value}

	resp, err := client.SendHTTPRequest(*method, tmp)
	if err != nil {
		log.Fatalf("sending http request failed with error %s", err)
	}

	log.Println(resp)

}
