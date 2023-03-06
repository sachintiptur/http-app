package main

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetKey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost:8080?key=Germany",
		httpmock.NewStringResponder(200, `[{"key": Germany, "value": "Berlin"}]`))

}
