package logging

import (
	"log"
	"net/http"
	"time"
)

// LogInfo Logging middle-ware handler struct
type LogInfo struct {
	handler http.Handler
}

// NewLogInfo middleware mux handler
func NewLogInfo(reqHandler http.Handler) *LogInfo {
	return &LogInfo{reqHandler}
}

// ServeHTTP Interface implementation for LogInfo
// Wraps the actual http handler functions with the
// logging middleware functions
func (l *LogInfo) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	log.Printf("METHOD: %s PATH: %s KEY: %s", req.Method, req.URL, req.URL.Query().Get("key"))
	l.handler.ServeHTTP(resp, req)
	log.Printf("Time elapsed: %v", time.Since(start))
}
