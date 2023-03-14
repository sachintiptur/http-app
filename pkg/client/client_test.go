package client

import (
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"go.uber.org/goleak"
)

var tmp Data

func TestCreateHTTPRequest(t *testing.T) {
	defer goleak.VerifyNone(t)
	var testcases = []struct {
		name               string
		method             string
		key                string
		value              string
		expectedResponse   string
		expectedStatusCode int
	}{
		{
			name:               "Test response for http PUT request",
			method:             http.MethodPut,
			key:                "Germany",
			value:              "Berlin",
			expectedResponse:   "Database updated with new key/value pair",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Test response for http PUT request",
			method:             http.MethodPut,
			key:                "Germany",
			value:              "Munich",
			expectedResponse:   "Updated the existing key/value pair",
			expectedStatusCode: http.StatusFound,
		},
		{
			name:               "Test response for http GET request",
			method:             http.MethodGet,
			key:                "Germany",
			expectedResponse:   "Data found for key Germany: Berlin",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Test response for http DELETE request",
			method:             http.MethodDelete,
			key:                "Germany",
			expectedResponse:   "Database entry deleted",
			expectedStatusCode: http.StatusOK,
		},
	}

	var url string

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			log.Println(tc.name)
			if tc.method == http.MethodPut {
				url = "http://localhost:8080?key=" + tc.key + "&value=" + tc.value
				tmp.Key = tc.key
				tmp.Val = tc.value
			} else {
				url = "http://localhost:8080" + "?key=" + tc.key
				tmp.Key = tc.key
			}

			httpmock.RegisterResponder(tc.method, url,
				httpmock.NewStringResponder(tc.expectedStatusCode, tc.expectedResponse))

			actualResponse, _ := SendHTTPRequest(tc.method, tmp)

			log.Println(actualResponse)
			require.Equal(t, tc.expectedResponse, actualResponse)

		})
	}

}
