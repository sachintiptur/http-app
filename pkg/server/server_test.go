package server

import (
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	database "github.com/sachintiptur/http-app/pkg/util"
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
			expectedStatus: http.StatusOK,
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

	var DatabaseUnderTest []DbStruct
	var url string
	// run tests for both file as database and local map as database
	DatabaseUnderTest = append(DatabaseUnderTest, DbStruct{&database.JsonData{}}, DbStruct{&database.MemData{}})

	for _, db := range DatabaseUnderTest {
		db.DbIntf.Init()
		log.Printf("Testing with %T as database\n", db.DbIntf)

		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.method == http.MethodPut {
					url = "http://localhost:8080?key=" + tc.key + "&value=" + tc.value
				} else {
					url = "http://localhost:8080" + "?key=" + tc.key

				}

				req := httptest.NewRequest(tc.method, url, nil)
				w := httptest.NewRecorder()
				db.ProcessHTTPRequests(w, req)

				require.Equal(t, tc.expectedStatus, w.Result().StatusCode)

			})
		}
	}

}
