package controller

import (
	"encoding/json"
	"net/http"
)

func Parse[T any](writer http.ResponseWriter, request *http.Request) (T, error) {
	input := new(T)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(input)

	if err != nil {
		return *input, err
	}

	return *input, nil
}
