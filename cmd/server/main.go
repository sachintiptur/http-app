package main

import (
	"flag"
	"log"
	"net/http"

	database "github.com/sachintiptur/http-app/pkg/database"
	"github.com/sachintiptur/http-app/pkg/middleware"
	"github.com/sachintiptur/http-app/pkg/server"
)

// HTTP server process
// Server address and db type can be passed as cli arguments
// Starts listening for http requests and calls the handler.
func main() {
	addr := flag.String("addr", ":8080", "Server address string")
	dbType := flag.String("db", "map", "Database to use, supported values are [map, json]")
	flag.Parse()

	if *dbType != database.MapType && *dbType != database.JsonType {
		flag.Usage()
		log.Fatal("Error: Unsupported database type")
	}
	// Map of supported database types
	db := map[string]database.Database{database.MapType: &database.InMemoryDatabase{}, database.JsonType: &database.JsonDatabase{}}
	var dbH server.DatabaseHandler

	// Initialise the database
	dbH.Db = db[*dbType]
	err := dbH.Db.Init()
	if err != nil {
		log.Fatalf("database initialisation failed: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", dbH.ProcessHTTPRequests)

	// Register the http handler
	mwMux := logging.NewLogInfo(mux)
	log.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(*addr, mwMux))
}
