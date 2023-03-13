package main

import (
	"flag"
	"fmt"
	"log"

	client "github.com/sachintiptur/http-app/pkg/client"
)

// Client CLI process main function
// takes http method, key and value as the command
// line arguments. Build and send appropriate HTTP request to
// the server.
func main() {

	method := flag.String("m", "", "HTTP method, supported methods are [GET, PUT, DELETE]")
	key := flag.String("k", "", "Key for the data")
	value := flag.String("v", "", "Value of the data")

	flag.Parse()

	var tmp = client.Data{Key: *key, Val: *value}

	resp, err := client.PrepareAndSendHTTPRequest(*method, tmp)
	if err != nil {
		fmt.Printf("PrepareAndSendHTTPRequest failed with error %s", err)
		return
	}

	log.Println(resp)

}
