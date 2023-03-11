package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	database "github.com/sachintiptur/http-app/util"
	"go.uber.org/goleak"
)

func TestProcessHTTPRequests(t *testing.T) {
	defer goleak.VerifyNone(t)
	var testcases = []struct {
		name           string
		method         string
		key            string
		value          string
		expectedStatus int
	}{

		{
			name:           "Test http PUT request ",
			method:         http.MethodPut,
			key:            "foo",
			value:          "bar",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test http GET request ",
			method:         http.MethodGet,
			key:            "foo",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test http patch request using PUT",
			method:         http.MethodPut,
			key:            "foo",
			value:          "barbar",
			expectedStatus: http.StatusFound,
		},
		{
			name:           "Test http DELETE request ",
			method:         http.MethodDelete,
			key:            "foo",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test http DELETE request for unknown entry ",
			method:         http.MethodDelete,
			key:            "foofoo",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Test http GET request failure scenario",
			method:         http.MethodGet,
			key:            "foo",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Test invalid key size in http PUT request ",
			method:         http.MethodPut,
			key:            "fooqourtbamflkgnfbhcbdhdjkt",
			value:          "bar",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Test invalid value size in http PUT request ",
			method:         http.MethodPut,
			key:            "foo",
			value:          "barqourtbamflkgnfbhcbdhdjktqourtbamflkgnfbhcbdhdjkt",
			expectedStatus: http.StatusBadRequest,
		},
	}

	var DatabaseUnderTest []dbStruct
	var url string
	// run tests for both file as database and local map as database
	DatabaseUnderTest = append(DatabaseUnderTest, dbStruct{&database.JsonData{}}, dbStruct{&database.MemData{}})

	for _, db := range DatabaseUnderTest {
		db.dbIntf.Init()
		log.Println()
		log.Printf("Testing with %T as database\n", db.dbIntf)
		for _, tc := range testcases {
			log.Println(tc.name)
			if tc.method == http.MethodPut {
				url = "http://localhost:8080?key=" + tc.key + "&value=" + tc.value
			} else {
				url = "http://localhost:8080" + "?key=" + tc.key

			}

			req := httptest.NewRequest(tc.method, url, nil)
			w := httptest.NewRecorder()
			db.processHTTPRequests(w, req)

			if tc.expectedStatus == w.Result().StatusCode {
				log.Println("PASSED")
			} else {
				log.Println("FAILED")
				log.Println(w.Result().Status)
			}

		}
	}

}
