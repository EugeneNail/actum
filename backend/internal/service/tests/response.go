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

func (response *Response) ReadData(reference any) *Response {
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(reference)
	Check(err)

	return response
}

func (response *Response) AssertStatus(status int) *Response {
	if response.StatusCode != status {
		response.t.Errorf("expected status %d, got %d", status, response.StatusCode)
		response.t.SkipNow()
	}

	return response
}

func (response *Response) AssertHasValidationErrors(fields []string) *Response {
	var errors map[string]string
	data, err := io.ReadAll(response.Body)
	Check(err)
	err = json.Unmarshal(data, &errors)
	Check(err)

	for _, field := range fields {
		if _, exists := errors[field]; !exists {
			response.t.Errorf(`expected validation error for field "%s" to be present`, field)
			response.t.SkipNow()
		}
	}

	return response
}
