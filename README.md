# HTTP client server


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
```

**Client**
```
Usage of ./client:
  -k string
    	Key for the data
  -m string
    	HTTP method
  -v string
    	Value of the data
```

## Example execution
**Server**
```
stiptur@mb02287 http-app % _build/server
2023/03/08 15:25:20 Server is listening...
2023/03/08 15:25:37 METHOD: PUT KEY: foo
2023/03/08 15:25:37 Time elapsed: 44.667µs
2023/03/08 15:25:50 METHOD: GET KEY: foo
2023/03/08 15:25:50 Time elapsed: 42µs
2023/03/08 15:25:59 METHOD: DELETE KEY: foo
2023/03/08 15:25:59 Time elapsed: 68.916µs
```

**Client**

```
stiptur@mb02287 http-app % ./_build/client -m=PUT -k=foo -v=bar
2023/03/08 15:25:37 200 OK
2023/03/08 15:25:37 Database updated with new key/value pair
stiptur@mb02287 http-app % ./_build/client -m=GET -k=foo       
2023/03/08 15:25:50 200 OK
2023/03/08 15:25:50 Data found for key foo: bar
stiptur@mb02287 http-app % ./_build/client -m=DELETE -k=foo
2023/03/08 15:25:59 200 OK
2023/03/08 15:25:59 Database entry deleted
stiptur@mb02287 http-app %
```

## Unit test coverage

**Server**
```
stiptur@mb02287 server % go tool cover -func cover.out  
github.com/sachintiptur/http-app/server/server.go:28:	InitDatabase		100.0%
github.com/sachintiptur/http-app/server/server.go:34:	ValidateData		85.7%
github.com/sachintiptur/http-app/server/server.go:52:	NewLogInfo		0.0%
github.com/sachintiptur/http-app/server/server.go:59:	ServeHTTP		0.0%
github.com/sachintiptur/http-app/server/server.go:67:	processGET		100.0%
github.com/sachintiptur/http-app/server/server.go:82:	processPUT		100.0%
github.com/sachintiptur/http-app/server/server.go:101:	processDELETE		100.0%
github.com/sachintiptur/http-app/server/server.go:118:	processHTTPRequests	90.0%
github.com/sachintiptur/http-app/server/server.go:143:	main			0.0%
total:							(statements)		72.7%
```

**Client**

```
stiptur@mb02287 client % go tool cover -func cover.out  
github.com/sachintiptur/http-app/client/client.go:17:	createHTTPRequest		92.3%
github.com/sachintiptur/http-app/client/client.go:47:	PrepareAndSendHTTPRequest	77.8%
github.com/sachintiptur/http-app/client/client.go:66:	main				0.0%
total:							(statements)			59.4%
```


