# HTTP client server
HTTP client server application using golang's net/http package.
This implements a basic http requests handling and it can be configured to
use either local map as a database or JSON file as a database.

## Build instructions

1. Compile server and client 
`make build`
2. Running unit tests 
`make test`
3. cleanup
`make clean`

## Usage

**Server**
```
Usage of ./server:
  -addr string
    	Server address string (default ":8080")
  -db string
    	Database to use, supported values are [map, json] (default "map")
```

**Client**
```
Usage of ./client:
  -k string
    	Key for the data
  -m string
    	HTTP method, supported methods are [GET, PUT, DELETE]
  -v string
    	Value of the data
```

## Example execution
**Server**
```
stiptur@mb02287 http-app % ./_build/server   
2023/03/12 22:13:10 Server is listening...
2023/03/12 22:13:25 METHOD: PUT PATH: /?key=foo-1&value=bar KEY: foo-1
2023/03/12 22:13:25 Time elapsed: 612.167µs
```

**Client**

```
stiptur@mb02287 http-app % ./_build/client -m=PUT -k=foo-1 -v=bar
2023/03/12 22:13:25 Database updated with new key/value pair

```

## Unit test coverage

**Server**
```
stiptur@mb02287 server % go test -coverprofile cover.out
2023/03/12 22:07:02 
2023/03/12 22:07:02 Testing with *database.JsonData as database
2023/03/12 22:07:02 Test http PUT request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http GET request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http patch request using PUT
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http DELETE request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http DELETE request for unknown entry 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http GET request failure scenario
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test invalid key size in http PUT request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test invalid value size in http PUT request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 
2023/03/12 22:07:02 Testing with *database.MemData as database
2023/03/12 22:07:02 Test http PUT request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http GET request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http patch request using PUT
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http DELETE request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http DELETE request for unknown entry 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test http GET request failure scenario
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test invalid key size in http PUT request 
2023/03/12 22:07:02 PASSED
2023/03/12 22:07:02 Test invalid value size in http PUT request 
2023/03/12 22:07:02 PASSED
PASS
coverage: 62.2% of statements
ok  	github.com/sachintiptur/http-app/server	0.249s
stiptur@mb02287 server % go tool cover -func cover.out  
github.com/sachintiptur/http-app/server/server.go:21:	ValidateData		100.0%
github.com/sachintiptur/http-app/server/server.go:37:	NewLogInfo		0.0%
github.com/sachintiptur/http-app/server/server.go:44:	ServeHTTP		0.0%
github.com/sachintiptur/http-app/server/server.go:52:	processGET		100.0%
github.com/sachintiptur/http-app/server/server.go:67:	processPUT		68.4%
github.com/sachintiptur/http-app/server/server.go:103:	processDELETE		75.0%
github.com/sachintiptur/http-app/server/server.go:131:	processHTTPRequests	90.0%
github.com/sachintiptur/http-app/server/server.go:156:	main			0.0%
total:							(statements)		62.2%
stiptur@mb02287 server % 

```

**Client**

```
stiptur@mb02287 client % go tool cover -func cover.out  
github.com/sachintiptur/http-app/client/client.go:17:	createHTTPRequest		92.3%
github.com/sachintiptur/http-app/client/client.go:47:	PrepareAndSendHTTPRequest	77.8%
github.com/sachintiptur/http-app/client/client.go:66:	main				0.0%
total:							(statements)			59.4%
```


