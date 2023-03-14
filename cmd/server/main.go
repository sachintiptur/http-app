package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/sachintiptur/http-app/pkg/server"
	database "github.com/sachintiptur/http-app/pkg/util"
)

// HTTP server process
// Server address and db type can be passed as cli arguments
// Starts listening for http requests and calls the handler.
func main() {
	addr := flag.String("addr", ":8080", "Server address string")
	dbType := flag.String("db", "map", "Database to use, supported values are [map, json]")
	flag.Parse()

	if *dbType != "map" && *dbType != "json" {
		flag.Usage()
		log.Fatal("Error: Unsupported database type")
	}
	// Map of supported database types
	var db = map[string]database.Database{"map": &database.MemData{}, "json": &database.JsonData{}}
	var dbS server.DbStruct

	// Initialise the database
	dbS.DbIntf = db[*dbType]
	err := dbS.DbIntf.Init()
	if err != nil {
		log.Fatalf("database initialisation failed: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", dbS.ProcessHTTPRequests)

	// Register the http handler
	mwMux := server.NewLogInfo(mux)
	log.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(*addr, mwMux))

}
