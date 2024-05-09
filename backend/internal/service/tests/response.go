package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type Response struct {
	*http.Response
	t *testing.T
}

func (response *Response) AssertStatus(status int) {
	if response.StatusCode != status {
		response.t.Errorf("expected status %d, got %d", status, response.StatusCode)
	}
}

func (response *Response) AssertHasValidationErrors(fields []string) {
	var errors map[string]string
	data, err := io.ReadAll(response.Body)
	Check(err)
	err = json.Unmarshal(data, &errors)
	Check(err)

	for _, field := range fields {
		if _, exists := errors[field]; !exists {
			response.t.Errorf(`expected validation error for field "%s" to be present`, field)
		}
	}
}

func (response *Response) AssertHasToken() {
	if !hasToken(response.Response) {
		response.t.Errorf("The response must have an Access-Token cookie")
	}
}

func (response *Response) AssertHasNoToken() {
	if hasToken(response.Response) {
		response.t.Errorf("The response must not have an Access-Token cookie")
	}
}

func hasToken(response *http.Response) bool {
	for _, cookie := range response.Cookies() {
		if cookie.Name == "Access-Token" && len(cookie.Value) > 0 {
			return true
		}
	}

	return false
}
